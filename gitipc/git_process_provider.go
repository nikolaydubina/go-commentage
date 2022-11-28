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

// ProcessGitProvider makes IPC calls to git to fech data.
// Parses human-readable git output.
// Current branch has to be branch of interest (eg, master).
type ProcessGitProvider struct {
	fileBlameDetails map[string]detailsForFile
}

func (s *ProcessGitProvider) processBlameFile(filename string) error {
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
	if s.fileBlameDetails == nil {
		s.fileBlameDetails = make(map[string]detailsForFile)
	}
	s.fileBlameDetails[filename] = detailsForFile{
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
	if _, ok := s.fileBlameDetails[filename]; !ok {
		if err := s.processBlameFile(filename); err != nil {
			return "", fmt.Errorf("can not process file: %w", err)
		}
	}
	commit, ok := s.fileBlameDetails[filename].CommitForLine[line]
	if !ok {
		return "", fmt.Errorf("line(%d) not found", line)
	}
	return commit, nil
}

func (s *ProcessGitProvider) LastUpdateForLine(filename string, line int) (lastUpdatedAt time.Time, err error) {
	if _, ok := s.fileBlameDetails[filename]; !ok {
		if err := s.processBlameFile(filename); err != nil {
			return lastUpdatedAt, fmt.Errorf("can not process file: %w", err)
		}
	}
	lastUpdatedAt, ok := s.fileBlameDetails[filename].LastUpdateForLine[line]
	if !ok {
		return lastUpdatedAt, fmt.Errorf("line(%d) not found", line)
	}
	return lastUpdatedAt, nil
}

func (s *ProcessGitProvider) NumCommitsFromRoot(commit string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", "--first-parent", commit)

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("can not execute command: %w", err)
	}

	if stderr.Len() > 0 {
		return 0, fmt.Errorf("stderr > 0: %s", stderr.String())
	}

	return strconv.Atoi(strings.TrimSpace(stdout.String()))
}
