package controllers

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/pkg/history"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"
	"github.com/go-openapi/runtime/middleware"
)

type ExecutionController struct {
	exec    *execution.Service
	history *history.Service
}

func NewExecutionController(exec *execution.Service, history *history.Service) (*ExecutionController, error) {
	if exec == nil {
		return nil, common.Error(common.ErrMissedArgument, "execution service")
	}

	if history == nil {
		return nil, common.Error(common.ErrMissedArgument, "history service")
	}

	return &ExecutionController{exec, history}, nil
}

func (ctl *ExecutionController) Create(params operations.CreateExecutionParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	ctx := context.Background()

	id, err := ctl.exec.Start(ctx, execution.Event{
		ProjectID: params.ProjectID,
		ScriptID:  *params.Body.ScriptID,
		CausedBy:  execution.CauseManual,
	})

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("script_id", *params.Body.ScriptID).
			Str("cause", execution.CauseManual.String()).
			Msg("failed to create a job")

		return http.InternalError()
	}

	logger.Info().
		Timestamp().
		Str("project_id", params.ProjectID).
		Str("script_id", *params.Body.ScriptID).
		Str("cause", execution.CauseManual.String()).
		Msg("create a job")

	return operations.NewCreateExecutionOK().WithPayload(id)
}

func (ctl *ExecutionController) Delete(params operations.DeleteExecutionParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	ctx := context.Background()

	err := ctl.exec.Cancel(ctx, params.ProjectID, params.JobID)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("job_id", params.JobID).
			Msg("failed to cancel a job")
	}

	logger.Error().
		Timestamp().
		Err(err).
		Str("project_id", params.ProjectID).
		Str("job_id", params.JobID).
		Msg("canceled a job")

	return operations.NewDeleteExecutionNoContent()
}

func (ctl *ExecutionController) Find(_ operations.FindExecutionsParams) middleware.Responder {
	return http.Bad("not implemented")
}

func (ctl *ExecutionController) Get(params operations.GetExecutionParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)

	found, err := ctl.history.Get(context.Background(), params.ProjectID, params.JobID)

	if err != nil {
		if err == common.ErrNotFound {
			return http.NotFound()
		}

		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Str("job_id", params.JobID).
			Msg("failed to find a job")

		return http.InternalError()
	}

	status := found.Status.String()
	cause := found.CausedBy.String()

	return operations.NewGetExecutionOK().WithPayload(&operations.GetExecutionOKBody{
		GetExecutionOKBodyAllOf0: operations.GetExecutionOKBodyAllOf0{
			JobID:     &found.JobID,
			ScriptID:  &found.ScriptID,
			ScriptRev: &found.ScriptRev,
			Status:    &status,
			Cause:     &cause,
		},
	})
}
