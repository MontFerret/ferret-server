package scripts

import "github.com/MontFerret/ferret-server/pkg/common/dal"

type (
	Identity struct {
		ProjectID string `json:"project_id"`
		ScriptID  string `json:"script_id"`
	}

	Persistence struct {
		Enabled bool `json:"enabled"`
	}

	Execution struct {
		Query  string                 `json:"query"`
		Params map[string]interface{} `json:"params"`
	}

	Script struct {
		dal.Definition
		Execution   Execution   `json:"execution"`
		Persistence Persistence `json:"persistence"`
	}

	UpdateScript struct {
		Script
		ID string `json:"id"`
	}

	ScriptEntity struct {
		dal.Entity
		Script
	}

	QueryResult struct {
		dal.QueryResult
		Data []ScriptEntity
	}
)
