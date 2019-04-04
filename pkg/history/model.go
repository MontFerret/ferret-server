package history

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"time"

	"github.com/MontFerret/ferret-server/pkg/execution"
)

type (
	// Represents a job history record
	// It gets updated during the job life cycle
	Record struct {
		JobID     string                 `json:"job_id"`
		ScriptID  string                 `json:"script_id"`
		ScriptRev string                 `json:"script_rev"`
		Params    map[string]interface{} `json:"params"`
		Status    execution.Status       `json:"status"`
		CausedBy  execution.Cause        `json:"cause"`
		StartedAt time.Time              `json:"started_at"`
		EndedAt   time.Time              `json:"ended_at"`
		Logs      []string               `json:"logs"`
		Error     error                  `json:"error"`
	}

	UpdateRecord struct {
		Record
		ID string `json:"id"`
	}

	RecordEntity struct {
		dal.Entity
		Record
	}

	QueryResult struct {
		dal.QueryResult
		Data []RecordEntity
	}
)
