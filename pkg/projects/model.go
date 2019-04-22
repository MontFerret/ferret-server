package projects

import "github.com/MontFerret/ferret-server/pkg/common/dal"

type (
	UpdateProject struct {
		Project
		ID string `json:"id"`
	}

	Project struct {
		dal.Definition
	}

	ProjectEntity struct {
		dal.Entity
		Project
	}

	QueryResult struct {
		dal.QueryResult
		Data []ProjectEntity
	}
)
