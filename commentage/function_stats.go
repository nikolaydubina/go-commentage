package commentage

import (
	"fmt"
	"time"
)

type FunctionStats struct {
	Name             string
	LastUpdatedAt    time.Time
	DocLastUpdatedAt time.Time
}

func (s FunctionStats) DocLastUpdatedBehindDays() time.Duration {
	return s.LastUpdatedAt.Sub(s.DocLastUpdatedAt)
}

func (s FunctionStats) String() string {
	return fmt.Sprintf(
		"%s last_updated_at(%v) doc_last_updated_at(%v) doc_last_updated_behind_days(%.2f)",
		s.Name,
		s.LastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedAt.Format(time.RFC3339),
		s.DocLastUpdatedBehindDays().Hours()/24,
	)
}
