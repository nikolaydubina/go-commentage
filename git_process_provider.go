package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ProcessGitProvider makes OS process calls to git to fech data.
// Utilizes human-readable git blame output.
// TODO: Caches data per file.
type ProcessGitProvider struct {
	Files map[string]DetailsForFile
}

func (s *ProcessGitProvider) processFile(filename string) error {
	cmd := exec.Command("git", "blame", "-t", "-e", filename)

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can not execute command: %w", err)
	}

	if stderr.Len() > 0 {
		return fmt.Errorf("stderr > 0: %s", stderr.String())
	}

	// parse git blame to get line details for whole file
	commitForLine := make(map[int]string)
	lastUpdateForLine := make(map[int]time.Time)

	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		// last line can be empty line in blame and in go file
		// not counted towards lines in source code
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		lineDetails, err := parseBlameLine(line)
		if err != nil {
			return fmt.Errorf("error at line(%s): %w", line, err)
		}

		commitForLine[lineDetails.lineNumber] = lineDetails.commit
		lastUpdateForLine[lineDetails.lineNumber] = lineDetails.createdAt
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("can not scan: %w", err)
	}

	// check that all lines collected
	if len(commitForLine) == 0 || len(lastUpdateForLine) == 0 {
		return errors.New("missing lines details")
	}
	if len(commitForLine) != len(lastUpdateForLine) {
		return fmt.Errorf("num lines with details mismatch")
	}
	// check line numbers are continious
	maxLine := 1
	for k := range commitForLine {
		if k > maxLine {
			maxLine = k
		}
	}
	for i := 1; i <= maxLine; i++ {
		if _, ok := commitForLine[i]; !ok {
			return fmt.Errorf("missing line(%d)", i)
		}
		if _, ok := lastUpdateForLine[i]; !ok {
			return fmt.Errorf("missing line(%d)", i)
		}
	}

	// store
	if s.Files == nil {
		s.Files = make(map[string]DetailsForFile)
	}
	s.Files[filename] = DetailsForFile{
		CommitForLine:     commitForLine,
		LastUpdateForLine: lastUpdateForLine,
	}

	return nil
}

type blameLine struct {
	lineNumber int
	commit     string
	createdAt  time.Time
}

// Example:
// ef0c9f0c5b8e pkg/util/node/node.go (<djmm@google.com>                 1464913558 -0700   2) Copyright 2015 The Kubernetes Authors.
func parseBlameLine(line string) (blameLine, error) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return blameLine{}, errors.New("wrong number of fields")
	}

	var lineDetails blameLine

	lineDetails.commit = fields[0]

	timestamp, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return blameLine{}, fmt.Errorf("non-interger timestamp(%s): %w", fields[3], err)
	}
	lineDetails.createdAt = time.Unix(timestamp, 0)

	lineNumberWithBracket := fields[5]
	if !strings.HasSuffix(lineNumberWithBracket, ")") {
		return blameLine{}, fmt.Errorf("expected to have '<line-number>)' but got(%s)", fields[5])
	}
	lineDetails.lineNumber, err = strconv.Atoi(lineNumberWithBracket[:len(lineNumberWithBracket)-1])
	if err != nil {
		return blameLine{}, fmt.Errorf("wrong number: %w", err)
	}

	return lineDetails, nil
}

func (s *ProcessGitProvider) CommitForLine(filename string, line int) (string, error) {
	if _, ok := s.Files[filename]; !ok {
		if err := s.processFile(filename); err != nil {
			return "", fmt.Errorf("can not process file: %w", err)
		}
	}
	commit, ok := s.Files[filename].CommitForLine[line]
	if !ok {
		return "", fmt.Errorf("line(%d) not found", line)
	}
	return commit, nil
}

func (s *ProcessGitProvider) LastUpdateForLine(filename string, line int) (time.Time, error) {
	if _, ok := s.Files[filename]; !ok {
		if err := s.processFile(filename); err != nil {
			return time.Time{}, fmt.Errorf("can not process file: %w", err)
		}
	}
	ts, ok := s.Files[filename].LastUpdateForLine[line]
	if !ok {
		return time.Time{}, fmt.Errorf("line(%d) not found", line)
	}
	return ts, nil
}
