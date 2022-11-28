package commentage

import (
	"fmt"
	"time"
)

type GitProider interface {
	CommitForLine(filename string, line int) (string, error)
	LastUpdateForLine(filename string, line int) (time.Time, error)
}

type SimpleStatsComputer struct {
	GitProider GitProider
}

func (s SimpleStatsComputer) lineRangeLastUpdatedAt(filename string, lineStart int, lineEnd int) (updatedAt time.Time, err error) {
	for i := lineStart; i <= lineEnd; i++ {
		lineUpdatedTime, err := s.GitProider.LastUpdateForLine(filename, i)
		if err != nil {
			return updatedAt, fmt.Errorf("can not get last updated_at for file(%s): %w", filename, err)
		}
		if lineUpdatedTime.After(updatedAt) {
			updatedAt = lineUpdatedTime
		}
	}
	return updatedAt, nil
}

func (s SimpleStatsComputer) CommmentLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	return s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
}

func (s SimpleStatsComputer) FunctionLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	return s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
}
