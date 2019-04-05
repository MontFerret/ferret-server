package dto

import (
	"github.com/MontFerret/ferret-server/server/http/api/models"

	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

func EntityFrom(from dal.Entity) models.Entity {
	return models.Entity{
		ID:       &from.ID,
		Rev:      &from.Rev,
		Metadata: *MetadataFrom(from.Metadata),
	}
}

func EntityTo(from models.Entity) (dal.Entity, error) {
	meta, err := MetadataTo(from.Metadata)

	if err != nil {
		return dal.Entity{}, err
	}

	return dal.Entity{
		ID:       *from.ID,
		Rev:      *from.Rev,
		Metadata: meta,
	}, nil
}
