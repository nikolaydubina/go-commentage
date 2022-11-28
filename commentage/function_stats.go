package commentage

import (
	"fmt"
	"strings"
	"time"
)

// FunctionStats is container for function and associated doc comment statistics.
type FunctionStats struct {
	Name string

	*TimeStats
	*CommitStats
}

func (s FunctionStats) String() string {
	var b strings.Builder
	b.WriteString(`"`)
	b.WriteString(s.Name)
	b.WriteString(`": `)

	if s.TimeStats != nil {
		b.WriteString(s.TimeStats.String())
	}

	if s.CommitStats != nil {
		b.WriteString(s.CommitStats.String())
	}

	return b.String()
}

func (s FunctionStats) StringVerbose() string {
	var b strings.Builder
	b.WriteString(`"`)
	b.WriteString(s.Name)
	b.WriteString(`": `)

	if s.TimeStats != nil {
		b.WriteString(s.TimeStats.StringVerbose())
	}

	if s.CommitStats != nil {
		b.WriteString(s.CommitStats.StringVerbose())
	}

	return b.String()
}

// CommitStats is commits related statistics of function and associated doc comment.
type CommitStats struct {
	LastCommit       string
	DocLastCommit    string
	DocBehindCommits int
}

func (s CommitStats) String() string {
	return fmt.Sprintf("doc_last_commit_behind(%d)", s.DocBehindCommits)
}

func (s CommitStats) StringVerbose() string {
	return fmt.Sprintf(
		"last_commit(%s) doc_last_commit(%s) doc_last_commit_behind(%d)",
		s.LastCommit,
		s.DocLastCommit,
		s.DocBehindCommits,
	)
}

// TimeStats is time related statistics of function and associated doc comment.
type TimeStats struct {
	LastUpdatedAt    time.Time
	DocLastUpdatedAt time.Time
}

func (s TimeStats) DocLastUpdatedBehindDays() time.Duration {
	return s.LastUpdatedAt.Sub(s.DocLastUpdatedAt)
}

func (s TimeStats) String() string {
	return fmt.Sprintf("doc_last_updated_behind_days(%.2f)", s.DocLastUpdatedBehindDays().Hours()/24)
}

func (s TimeStats) StringVerbose() string {
	return fmt.Sprintf(
		"last_updated_at(%v) doc_last_updated_at(%v) doc_last_updated_behind_days(%.2f)",
		s.LastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedBehindDays().Hours()/24,
	)
}
