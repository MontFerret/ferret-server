package dto

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/server/http/api/models"
	"github.com/pkg/errors"
	"time"
)

func MetadataFrom(from dal.Metadata) *models.Metadata {
	var createdAt string
	var updatedAt string

	if from.CreatedAt != nil && from.CreatedAt.IsZero() == false {
		createdAt = from.CreatedAt.String()
	}

	if from.UpdateAt != nil && from.UpdateAt.IsZero() == false {
		updatedAt = from.UpdateAt.String()
	}

	return &models.Metadata{
		CreatedAt: &createdAt,
		UpdatedAt: updatedAt,
	}
}

func MetadataTo(from models.Metadata) (dal.Metadata, error) {
	var createdAt *time.Time
	var updatedAt *time.Time

	if from.CreatedAt != nil {
		str := *from.CreatedAt

		parsed, err := time.Parse(time.RFC3339, str)

		if err != nil {
			return dal.Metadata{}, errors.Wrap(err, "parse created_at")
		}

		createdAt = &parsed
	} else {
		ts := time.Now()
		// fallback
		createdAt = &ts
	}

	if from.UpdatedAt != "" {
		parsed, err := time.Parse(time.RFC3339, from.UpdatedAt)

		if err != nil {
			return dal.Metadata{}, errors.Wrap(err, "parse updated_at")
		}

		updatedAt = &parsed
	}

	return dal.Metadata{
		CreatedAt: createdAt,
		UpdateAt:  updatedAt,
	}, nil
}
