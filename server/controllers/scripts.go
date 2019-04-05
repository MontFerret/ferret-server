package controllers

import (
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/models"
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
		Definition:  dto.DefinitionTo(params.Body.Definition),
		Execution:   *ctl.fromScriptExecutionDto(params.Body.Execution),
		Persistence: *ctl.fromScriptPersistenceDto(params.Body.Persistence),
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

	e := dto.EntityFrom(out)
	return operations.NewCreateScriptCreated().WithPayload(&e)
}

func (ctl *Scripts) UpdateScript(params operations.UpdateScriptParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	script := scripts.UpdateScript{
		Script: scripts.Script{
			Definition:  dto.DefinitionTo(params.Body.Definition),
			Execution:   *ctl.fromScriptExecutionDto(params.Body.Execution),
			Persistence: *ctl.fromScriptPersistenceDto(params.Body.Persistence),
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

	e := dto.EntityFrom(out)

	return operations.NewUpdateScriptOK().WithPayload(&e)
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

	return operations.NewGetScriptOK().WithPayload(&models.ScriptOutputDetailed{
		ScriptEntity: models.ScriptEntity{
			Entity:       dto.EntityFrom(out.Entity),
			ScriptCommon: ctl.toScriptCommonDto(out),
		},
	})
}

func (ctl *Scripts) FindScripts(params operations.FindScriptsParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	query := dal.Query{
		Pagination: dto.PaginationTo(params.Count, params.Cursor),
	}

	out, err := ctl.service.FindScripts(params.HTTPRequest.Context(), params.ProjectID, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Uint64("count", query.Pagination.Count).
			Str("cursor", query.Pagination.Cursor.String()).
			Msg("failed to find scripts")

		return http.InternalError()
	}

	data := make([]*models.ScriptOutput, 0, len(out.Data))

	for _, s := range out.Data {
		script := s

		data = append(data, &models.ScriptOutput{
			Entity:     dto.EntityFrom(script.Entity),
			Definition: dto.DefinitionFrom(script.Definition),
		})
	}

	return operations.NewFindScriptsOK().WithPayload(&operations.FindScriptsOKBody{
		Data:         data,
		SearchResult: dto.SearchResultFrom(out.QueryResult),
	})
}

func (ctl *Scripts) toScriptCommonDto(script scripts.ScriptEntity) models.ScriptCommon {
	return models.ScriptCommon{
		Definition: dto.DefinitionFrom(script.Definition),
		Persistence: &models.ScriptPersistence{
			Enabled: &script.Persistence.Enabled,
		},
		Execution: &models.ScriptExecution{
			Query:  &script.Execution.Query,
			Params: dto.ExecutionParamsFrom(script.Execution.Params),
		},
	}
}

func (ctl *Scripts) fromScriptExecutionDto(exec *models.ScriptExecution) *scripts.Execution {
	return &scripts.Execution{
		Query:  *exec.Query,
		Params: dto.ExecutionParamsTo(exec.Params),
	}
}

func (ctl *Scripts) fromScriptPersistenceDto(exec *models.ScriptPersistence) *scripts.Persistence {
	return &scripts.Persistence{
		Enabled: *exec.Enabled,
	}
}
