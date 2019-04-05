package dto

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/server/http/api/models"
)

func DefinitionFrom(from dal.Definition) models.Definition {
	return models.Definition{
		Name:        &from.Name,
		Description: from.Description,
	}
}

func DefinitionTo(from models.Definition) dal.Definition {
	return dal.Definition{
		Name:        *from.Name,
		Description: from.Description,
	}
}
