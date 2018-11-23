package controllers

import (
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"
	"github.com/go-openapi/runtime/middleware"
)

type ProjectsController struct {
	service *projects.Service
}

func NewProjectsController(service *projects.Service) (*ProjectsController, error) {
	if service == nil {
		return nil, common.Error(common.ErrMissedArgument, "exec")
	}

	return &ProjectsController{service}, nil
}

func (ctl *ProjectsController) CreateProject(params operations.CreateProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.CreateProject(params.HTTPRequest.Context(), projects.Project{
		Name:        *params.Body.Name,
		Description: params.Body.Description,
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

	createdAt, _ := dto.ToMetadataDates(out.Metadata)

	return operations.NewCreateProjectCreated().WithPayload(&operations.CreateProjectCreatedBody{
		ID:        &out.ID,
		Rev:       &out.Rev,
		CreatedAt: &createdAt,
	})
}

func (ctl *ProjectsController) UpdateProject(params operations.UpdateProjectParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	out, err := ctl.service.UpdateProject(params.HTTPRequest.Context(), projects.UpdateProject{
		Project: projects.Project{
			Name:        *params.Body.Name,
			Description: params.Body.Description,
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

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewUpdateProjectOK().WithPayload(&operations.UpdateProjectOKBody{
		ID:        &out.ID,
		Rev:       &out.Rev,
		CreatedAt: &createdAt,
		UpdatedAt: updatedAt,
	})
}

func (ctl *ProjectsController) DeleteProject(params operations.DeleteProjectParams) middleware.Responder {
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

func (ctl *ProjectsController) GetProject(params operations.GetProjectParams) middleware.Responder {
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

	createdAt, updatedAt := dto.ToMetadataDates(out.Metadata)

	return operations.NewGetProjectOK().WithPayload(&operations.GetProjectOKBody{
		Name:        &out.Name,
		Description: out.Description,
		GetProjectOKBodyAllOf0: operations.GetProjectOKBodyAllOf0{
			ID:        &out.ID,
			Rev:       &out.Rev,
			CreatedAt: &createdAt,
			UpdatedAt: updatedAt,
		},
	})
}

func (ctl *ProjectsController) FindProjects(params operations.FindProjectsParams) middleware.Responder {
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

	out, err := ctl.service.FindProjects(params.HTTPRequest.Context(), query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Uint("page", query.Pagination.Page).
			Uint("size", query.Pagination.Size).
			Msg("failed to find projects")

		return http.InternalError()
	}

	res := make([]*operations.FindProjectsOKBodyItems0, 0, len(out))

	for _, p := range out {
		createdAt, updatedAt := dto.ToMetadataDates(p.Metadata)

		res = append(res, &operations.FindProjectsOKBodyItems0{
			Name: &p.Name,
			FindProjectsOKBodyItems0AllOf0: operations.FindProjectsOKBodyItems0AllOf0{
				ID:        &p.ID,
				Rev:       &p.Rev,
				CreatedAt: &createdAt,
				UpdatedAt: updatedAt,
			},
		})
	}

	return operations.NewFindProjectsOK().WithPayload(res)
}
