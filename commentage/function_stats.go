package commentage

import (
	"fmt"
	"time"
)

type FunctionStats struct {
	Name             string
	LastCommit       string
	LastUpdatedAt    time.Time
	DocLastCommit    string
	DocLastUpdatedAt time.Time
	DocBehindCommits int
}

func (s FunctionStats) DocLastUpdatedBehindDays() time.Duration {
	return s.LastUpdatedAt.Sub(s.DocLastUpdatedAt)
}

func (s FunctionStats) String() string {
	return fmt.Sprintf(
		"\"%s\": last_updated_at(%v) doc_last_updated_at(%v) doc_last_updated_behind_days(%.2f) last_commit(%s) doc_last_commit(%s) doc_last_commit_behind(%d)",
		s.Name,
		s.LastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedBehindDays().Hours()/24,
		s.LastCommit,
		s.DocLastCommit,
		s.DocBehindCommits,
	)
}
