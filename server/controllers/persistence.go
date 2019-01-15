package controllers

import (
	"context"
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

type PersistenceController struct {
	service *persistence.Service
}

func NewPersistenceController(service *persistence.Service) (*PersistenceController, error) {
	if service == nil {
		return nil, common.Error(common.ErrMissedArgument, "persistence service")
	}

	return &PersistenceController{service}, nil
}

func (ctl *PersistenceController) FindAll(params operations.FindProjectDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	var size uint = 10
	var page uint = 1

	if params.Size != nil {
		size = uint(*params.Size)
		page = uint(*params.Page)
	}

	query := dal.Query{
		Pagination: dal.Pagination{
			Size: size,
			Page: page,
		},
	}

	ctx := context.Background()
	out, err := ctl.service.FindProjectRecords(ctx, params.ProjectID, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Msg("failed to find project data")

		return http.InternalError()
	}

	payload := make([]*operations.FindProjectDataOKBodyItems0, 0, len(out))

	for _, i := range out {
		el := i
		createdAt, updatedAt := dto.ToMetadataDates(el.Metadata)

		payload = append(payload, &operations.FindProjectDataOKBodyItems0{
			FindProjectDataOKBodyItems0AllOf0: operations.FindProjectDataOKBodyItems0AllOf0{
				ID:        &el.ID,
				Rev:       &el.Rev,
				CreatedAt: &createdAt,
				UpdatedAt: updatedAt,
			},
			JobID:     &el.JobID,
			ScriptID:  &el.ScriptID,
			ScriptRev: el.ScriptRev,
		})
	}

	return operations.NewFindProjectDataOK().WithPayload(payload)
}

func (ctl *PersistenceController) Find(params operations.FindScriptDataParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	var size uint = 10
	var page uint = 1

	if params.Size != nil {
		size = uint(*params.Size)
		page = uint(*params.Page)
	}

	query := dal.Query{
		Pagination: dal.Pagination{
			Size: size,
			Page: page,
		},
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
			Msg("failed to find script data")

		return http.InternalError()
	}

	payload := make([]*operations.FindScriptDataOKBodyItems0, 0, len(out))

	for _, i := range out {
		el := i
		createdAt, updatedAt := dto.ToMetadataDates(el.Metadata)

		payload = append(payload, &operations.FindScriptDataOKBodyItems0{
			FindScriptDataOKBodyItems0AllOf0: operations.FindScriptDataOKBodyItems0AllOf0{
				ID:        &el.ID,
				Rev:       &el.Rev,
				CreatedAt: &createdAt,
				UpdatedAt: updatedAt,
			},
			JobID:     &el.JobID,
			ScriptID:  &el.ScriptID,
			ScriptRev: el.ScriptRev,
		})
	}

	return operations.NewFindScriptDataOK().WithPayload(payload)
}

func (ctl *PersistenceController) Get(params operations.GetScriptDataParams) middleware.Responder {
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

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewGetScriptDataOK().WithPayload(&operations.GetScriptDataOKBody{
		GetScriptDataOKBodyAllOf0: operations.GetScriptDataOKBodyAllOf0{
			GetScriptDataOKBodyAllOf0AllOf0: operations.GetScriptDataOKBodyAllOf0AllOf0{
				ID:        &out.ID,
				Rev:       &out.Rev,
				CreatedAt: &createdAt,
				UpdatedAt: updatedAt,
			},
			ScriptID:  &out.ScriptID,
			ScriptRev: &out.ScriptRev,
			JobID:     &out.JobID,
			Value:     out.Data,
		},
	})
}

func (ctl *PersistenceController) Update(params operations.UpdateScriptDataParams) middleware.Responder {
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

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewUpdateScriptDataOK().WithPayload(&operations.UpdateScriptDataOKBody{
		ID:        &out.ID,
		Rev:       &out.Rev,
		CreatedAt: &createdAt,
		UpdatedAt: updatedAt,
	})
}

func (ctl *PersistenceController) Delete(params operations.DeleteScriptDataParams) middleware.Responder {
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
