package gitipc

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

type detailsForFile struct {
	CommitForLine     map[int]string
	LastUpdateForLine map[int]time.Time
}

// ProcessGitProvider makes OS process calls to git to fech data.
// Utilizes human-readable git blame output.
// TODO: Caches data per file.
type ProcessGitProvider struct {
	Files map[string]detailsForFile
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
	lineNumber := 1 // start from 1 to match go ast
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

		commitForLine[lineNumber] = lineDetails.commit
		lastUpdateForLine[lineNumber] = lineDetails.createdAt

		lineNumber++
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

	// store
	if s.Files == nil {
		s.Files = make(map[string]detailsForFile)
	}
	s.Files[filename] = detailsForFile{
		CommitForLine:     commitForLine,
		LastUpdateForLine: lastUpdateForLine,
	}

	return nil
}

type lineDetails struct {
	commit    string
	createdAt time.Time
}

func parseBlameLine(line string) (ld lineDetails, err error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return ld, errors.New("wrong number of fields")
	}

	ld.commit = fields[0]

	// created at
	idxTimeStamp := 3

	// filename sometimes can be omitted
	if fields[1][0] == '(' {
		idxTimeStamp = 2
	}

	rawTimeStamp := fields[idxTimeStamp]
	timestamp, err := strconv.ParseInt(rawTimeStamp, 10, 64)
	if err != nil {
		return ld, fmt.Errorf("bad timestamp(%s): %w", rawTimeStamp, err)
	}
	ld.createdAt = time.Unix(timestamp, 0)

	return ld, err
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
