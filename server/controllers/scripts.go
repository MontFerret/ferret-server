package controllers

import (
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"

	"github.com/go-openapi/runtime/middleware"
)

type Scripts struct {
	service *scripts.Service
}

func NewScripts(service *scripts.Service) (*Scripts, error) {
	if service == nil {
		return nil, common.Error(common.ErrMissedArgument, "exec")
	}

	return &Scripts{service}, nil
}

func (ctl *Scripts) CreateScript(params operations.CreateScriptParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	entity := scripts.Script{
		Name:        *params.Body.Name,
		Description: params.Body.Description,
		Execution: scripts.Execution{
			Query:  *params.Body.Execution.Query,
			Params: params.Body.Execution.Params,
		},
		Persistence: scripts.Persistence{
			Enabled: *params.Body.Persistence.Enabled,
		},
	}

	out, err := ctl.service.CreateScript(params.HTTPRequest.Context(), params.ProjectID, entity)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("name", *params.Body.Name).
			Msg("failed to create new script")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("project_id", params.ProjectID).
		Str("name", *params.Body.Name).
		Msg("created new script")

	createdAt, _ := dto.ToMetadataDates(out.Metadata)

	return operations.NewCreateScriptCreated().WithPayload(&operations.CreateScriptCreatedBody{
		ID:        &out.ID,
		Rev:       &out.Rev,
		CreatedAt: &createdAt,
	})
}

func (ctl *Scripts) UpdateScript(params operations.UpdateScriptParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	script := scripts.UpdateScript{
		Script: scripts.Script{
			Name:        *params.Body.Name,
			Description: params.Body.Description,
			Execution: scripts.Execution{
				Query:  *params.Body.Execution.Query,
				Params: params.Body.Execution.Params,
			},
			Persistence: scripts.Persistence{
				Enabled: *params.Body.Persistence.Enabled,
			},
		},
		ID: params.ScriptID,
	}

	out, err := ctl.service.UpdateScript(params.HTTPRequest.Context(), params.ProjectID, script)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("id", params.ScriptID).
			Str("name", *params.Body.Name).
			Msg("failed to create new script")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("project_id", params.ProjectID).
		Str("id", params.ScriptID).
		Str("name", *params.Body.Name).
		Msg("updated script")

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewUpdateScriptOK().WithPayload(&operations.UpdateScriptOKBody{
		ID:        &out.ID,
		Rev:       &out.Rev,
		CreatedAt: &createdAt,
		UpdatedAt: updatedAt,
	})
}

func (ctl *Scripts) DeleteScript(params operations.DeleteScriptParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	err := ctl.service.DeleteScript(params.HTTPRequest.Context(), params.ProjectID, params.ScriptID)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("id", params.ScriptID).
			Msg("failed to delete script")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("project_id", params.ProjectID).
		Str("id", params.ScriptID).
		Msg("deleted script")

	return operations.NewDeleteScriptNoContent()
}

func (ctl *Scripts) GetScript(params operations.GetScriptParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.GetScript(params.HTTPRequest.Context(), params.ProjectID, params.ScriptID)

	if err != nil {
		if err == common.ErrNotFound {
			return http.NotFound()
		}

		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("id", params.ScriptID).
			Msg("failed to get script")

		return http.InternalError()
	}

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewGetScriptOK().WithPayload(&operations.GetScriptOKBody{
		GetScriptOKBodyAllOf0: operations.GetScriptOKBodyAllOf0{
			ID:        &out.ID,
			Rev:       &out.Rev,
			CreatedAt: &createdAt,
			UpdatedAt: updatedAt,
		},
		Name:        &out.Name,
		Description: out.Description,
		Execution: &operations.GetScriptOKBodyAO1Execution{
			Query:  &out.Execution.Query,
			Params: out.Execution.Params,
		},
		Persistence: &operations.GetScriptOKBodyAO1Persistence{
			Enabled: &out.Persistence.Enabled,
		},
	})
}

func (ctl *Scripts) FindScripts(params operations.FindScriptsParams) middleware.Responder {
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

	out, err := ctl.service.FindScripts(params.HTTPRequest.Context(), params.ProjectID, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Uint("page", query.Pagination.Page).
			Uint("size", query.Pagination.Size).
			Msg("failed to find scripts")

		return http.InternalError()
	}

	res := make([]*operations.FindScriptsOKBodyItems0, 0, len(out))

	for _, i := range out {
		p := i
		createdAt, updatedAt := dto.ToMetadataDates(p.Metadata)

		res = append(res, &operations.FindScriptsOKBodyItems0{
			FindScriptsOKBodyItems0AllOf0: operations.FindScriptsOKBodyItems0AllOf0{
				ID:        &p.ID,
				Rev:       &p.Rev,
				CreatedAt: &createdAt,
				UpdatedAt: updatedAt,
			},
			Name:        &p.Name,
			Description: p.Description,
		})
	}

	return operations.NewFindScriptsOK().WithPayload(res)
}
