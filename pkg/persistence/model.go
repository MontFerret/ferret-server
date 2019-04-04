package persistence

import "github.com/MontFerret/ferret-server/pkg/common/dal"

type (
	Record struct {
		JobID     string      `json:"job_id"`
		ScriptID  string      `json:"script_id"`
		ScriptRev string      `json:"script_rev"`
		Data      interface{} `json:"data"`
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
