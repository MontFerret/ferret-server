package dto

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

func ToMetadataDates(from dal.Metadata) (string, string) {
	var createdAt string
	var updatedAt string

	if from.CreatedAt.IsZero() == false {
		createdAt = from.CreatedAt.String()
	}

	if from.UpdateAt.IsZero() == false {
		updatedAt = from.UpdateAt.String()
	}

	return createdAt, updatedAt
}
