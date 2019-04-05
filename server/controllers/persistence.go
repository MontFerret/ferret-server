package controllers

import (
	"context"
	"github.com/MontFerret/ferret-server/server/http/api/models"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/persistence"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"

	"github.com/go-openapi/runtime/middleware"
)

type Persistence struct {
	service *persistence.Service
}

func NewPersistence(service *persistence.Service) (*Persistence, error) {
	if service == nil {
		return nil, common.Error(common.ErrMissedArgument, "persistence service")
	}

	return &Persistence{service}, nil
}

func (ctl *Persistence) FindAll(params operations.FindProjectDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	query := dal.Query{
		Pagination: dto.PaginationTo(params.Count, params.Cursor),
	}

	ctx := context.Background()
	out, err := ctl.service.FindProjectRecords(ctx, params.ProjectID, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Uint64("count", query.Pagination.Count).
			Str("cursor", query.Pagination.Cursor.String()).
			Msg("failed to find project data")

		return http.InternalError()
	}

	data := make([]*models.DataOutput, 0, len(out.Data))

	for _, i := range out.Data {
		data = append(data, ctl.toDataOutputDto(i))
	}

	return operations.NewFindProjectDataOK().WithPayload(&operations.FindProjectDataOKBody{
		Data:         data,
		SearchResult: dto.SearchResultFrom(out.QueryResult),
	})
}

func (ctl *Persistence) Find(params operations.FindScriptDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	query := dal.Query{
		Pagination: dto.PaginationTo(params.Count, params.Cursor),
	}

	ctx := context.Background()
	out, err := ctl.service.FindScriptRecords(ctx, scripts.Identity{
		ProjectID: params.ProjectID,
		ScriptID:  params.ScriptID,
	}, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("script_id", params.ScriptID).
			Uint64("count", query.Pagination.Count).
			Str("cursor", query.Pagination.Cursor.String()).
			Msg("failed to find script data")

		return http.InternalError()
	}

	data := make([]*models.DataOutput, 0, len(out.Data))

	for _, i := range out.Data {
		data = append(data, ctl.toDataOutputDto(i))
	}

	return operations.NewFindScriptDataOK().WithPayload(&operations.FindScriptDataOKBody{
		Data:         data,
		SearchResult: dto.SearchResultFrom(out.QueryResult),
	})
}

func (ctl *Persistence) Get(params operations.GetScriptDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	ctx := context.Background()

	out, err := ctl.service.GetRecord(ctx, scripts.Identity{
		ProjectID: params.ProjectID,
		ScriptID:  params.ScriptID,
	}, params.DataID)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("script_id", params.ScriptID).
			Msg("failed to get a script data")

		return http.InternalError()
	}

	return operations.NewGetScriptDataOK().WithPayload(&models.DataOutputDetailed{
		DataEntity: models.DataEntity{
			Entity: dto.EntityFrom(out.Entity),
			DataCommon: models.DataCommon{
				ScriptID:  &out.ScriptID,
				ScriptRev: &out.ScriptRev,
				JobID:     &out.JobID,
				Value:     out.Data,
			},
		},
	})
}

func (ctl *Persistence) Update(params operations.UpdateScriptDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	ctx := context.Background()

	out, err := ctl.service.UpdateRecord(ctx, params.ProjectID, persistence.Record{
		Data: params.Body.Value,
	})

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("script_id", params.ScriptID).
			Str("id", params.DataID).
			Msg("failed to update a script data")

		return http.InternalError()
	}

	e := dto.EntityFrom(out)
	return operations.NewUpdateScriptDataOK().WithPayload(&e)
}

func (ctl *Persistence) Delete(params operations.DeleteScriptDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	ctx := context.Background()

	err := ctl.service.DeleteRecord(ctx, scripts.Identity{
		ProjectID: params.ProjectID,
		ScriptID:  params.ScriptID,
	}, params.DataID)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("script_id", params.ScriptID).
			Str("id", params.DataID).
			Msg("failed to update a script data")

		return http.InternalError()
	}

	return operations.NewDeleteScriptDataNoContent()
}

func (ctl *Persistence) toDataOutputDto(data persistence.RecordEntity) *models.DataOutput {
	return &models.DataOutput{
		Entity:    dto.EntityFrom(data.Entity),
		JobID:     &data.JobID,
		ScriptID:  &data.ScriptID,
		ScriptRev: data.ScriptRev,
	}
}
