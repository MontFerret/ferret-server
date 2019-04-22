package controllers

import (
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/models"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"

	"github.com/go-openapi/runtime/middleware"
)

type Projects struct {
	service *projects.Service
}

func NewProjects(service *projects.Service) (*Projects, error) {
	if service == nil {
		return nil, common.Error(common.ErrMissedArgument, "exec")
	}

	return &Projects{service}, nil
}

func (ctl *Projects) CreateProject(params operations.CreateProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.CreateProject(params.HTTPRequest.Context(), projects.Project{
		Definition: dto.DefinitionTo(params.Body.Definition),
	})

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("name", *params.Body.Name).
			Msgf("failed to create new project")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("id", out.ID).
		Str("name", *params.Body.Name).
		Msg("created new project")

	e := dto.EntityFrom(out)

	return operations.NewCreateProjectCreated().WithPayload(&e)
}

func (ctl *Projects) UpdateProject(params operations.UpdateProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.UpdateProject(params.HTTPRequest.Context(), projects.UpdateProject{
		Project: projects.Project{
			Definition: dto.DefinitionTo(params.Body.Definition),
		},
		ID: params.ProjectID,
	})

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("id", params.ProjectID).
			Str("name", *params.Body.Name).
			Msg("failed to update project")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("id", out.ID).
		Str("name", *params.Body.Name).
		Msg("updated project")

	e := dto.EntityFrom(out)
	return operations.NewUpdateProjectOK().WithPayload(&e)
}

func (ctl *Projects) DeleteProject(params operations.DeleteProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	err := ctl.service.DeleteProject(params.HTTPRequest.Context(), params.ProjectID)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("id", params.ProjectID).
			Msg("failed to delete project")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("id", params.ProjectID).
		Msg("deleted project")

	return operations.NewDeleteProjectNoContent()
}

func (ctl *Projects) GetProject(params operations.GetProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.GetProject(params.HTTPRequest.Context(), params.ProjectID)

	if err != nil {
		if err == common.ErrNotFound {
			return http.NotFound()
		}

		logger.Error().
			Timestamp().
			Err(err).
			Str("id", params.ProjectID).
			Msg("failed to get project")

		return http.InternalError()
	}

	entity := dto.EntityFrom(out.Entity)

	return operations.NewGetProjectOK().WithPayload(&models.ProjectOutputDetailed{
		ProjectEntity: models.ProjectEntity{
			Entity: entity,
			ProjectCommon: models.ProjectCommon{
				Definition: dto.DefinitionFrom(out.Definition),
			},
		},
	})
}

func (ctl *Projects) FindProjects(params operations.FindProjectsParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	query := dal.Query{
		Pagination: dto.PaginationTo(params.Count, params.Cursor),
	}

	out, err := ctl.service.FindProjects(params.HTTPRequest.Context(), query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Uint64("count", query.Pagination.Count).
			Str("cursor", query.Pagination.Cursor.String()).
			Msg("failed to find projects")

		return http.InternalError()
	}

	res := make([]*models.ProjectOutput, 0, len(out.Data))

	for _, p := range out.Data {
		project := p

		res = append(res, &models.ProjectOutput{
			Entity:     dto.EntityFrom(project.Entity),
			Definition: dto.DefinitionFrom(project.Definition),
		})
	}

	return operations.NewFindProjectsOK().WithPayload(&operations.FindProjectsOKBody{
		Data:         res,
		SearchResult: dto.SearchResultFrom(out.QueryResult),
	})
}
