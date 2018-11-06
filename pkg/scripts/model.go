package scripts

import "github.com/MontFerret/ferret-server/pkg/common/dal"

type (
	Persistence struct {
		Local  string   `json:"local"`
		Remote []string `json:"remote"`
	}

	Execution struct {
		Query  string                 `json:"content"`
		Params map[string]interface{} `json:"params"`
	}

	Script struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
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
)
