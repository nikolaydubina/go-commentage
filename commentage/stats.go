package commentage

import (
	"fmt"
	"time"
)

type GitProider interface {
	CommitForLine(filename string, line int) (string, error)
	LastUpdateForLine(filename string, line int) (time.Time, error)
	NumCommitsFromRoot(commit string) (int, error)
}

type SimpleStatsComputer struct {
	GitProider GitProider
}

func (s SimpleStatsComputer) lineRangeLastUpdatedAt(filename string, lineStart int, lineEnd int) (lineNumber int, updatedAt time.Time, err error) {
	for i := lineStart; i <= lineEnd; i++ {
		lineUpdatedTime, err := s.GitProider.LastUpdateForLine(filename, i)
		if err != nil {
			return lineNumber, updatedAt, fmt.Errorf("can not get last updated_at for file(%s): %w", filename, err)
		}
		if lineUpdatedTime.After(updatedAt) {
			lineNumber = i
			updatedAt = lineUpdatedTime
		}
	}
	return lineNumber, updatedAt, nil
}

func (s SimpleStatsComputer) CommmentLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	_, t, err := s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
	return t, err
}

func (s SimpleStatsComputer) FunctionLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	_, t, err := s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
	return t, err
}

// LastCommitForRange selects commit for last updated line.
// Not using true ordering of commits, since it is computationally prohibitive.
func (s SimpleStatsComputer) LastCommitForRange(filename string, lineStart int, lineEnd int) (string, error) {
	ln, _, err := s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
	if err != nil {
		return "", fmt.Errorf("can not get last updated line for file(%s): %w", filename, err)
	}
	return s.GitProider.CommitForLine(filename, ln)
}

func (s SimpleStatsComputer) CommitDifference(filename string, fromLineStart, fromLineEnd, toLineStart, toLineEnd int) (int, error) {
	fromMaxCommit, err := s.LastCommitForRange(filename, fromLineStart, fromLineEnd)
	if err != nil {
		return 0, fmt.Errorf("can not get last commit for range(%s:%d:%d): %w", filename, fromLineStart, fromLineEnd, err)
	}

	toMaxCommit, err := s.LastCommitForRange(filename, toLineStart, toLineEnd)
	if err != nil {
		return 0, fmt.Errorf("can not get last commit for range(%s:%d:%d): %w", filename, toLineStart, toLineEnd, err)
	}

	fromNumMaxCommitsFromRoot, err := s.GitProider.NumCommitsFromRoot(fromMaxCommit)
	if err != nil {
		return 0, fmt.Errorf("can not get num commits from root for commit(%s): %w", fromMaxCommit, err)
	}

	toNumMaxCommitsFromRoot, err := s.GitProider.NumCommitsFromRoot(toMaxCommit)
	if err != nil {
		return 0, fmt.Errorf("can not get num commits from root for commit(%s): %w", toMaxCommit, err)
	}

	return toNumMaxCommitsFromRoot - fromNumMaxCommitsFromRoot, nil
}
