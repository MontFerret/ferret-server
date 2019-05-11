package controllers

import (
	"context"
	"github.com/MontFerret/ferret-server/server/controllers/dto"
	"github.com/MontFerret/ferret-server/server/http/api/models"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/pkg/history"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"

	"github.com/go-openapi/runtime/middleware"
)

type Execution struct {
	exec    *execution.Service
	history *history.Service
}

func NewExecution(exec *execution.Service, history *history.Service) (*Execution, error) {
	if exec == nil {
		return nil, common.Error(common.ErrMissedArgument, "execution service")
	}

	if history == nil {
		return nil, common.Error(common.ErrMissedArgument, "history service")
	}

	return &Execution{exec, history}, nil
}

func (ctl *Execution) Create(params operations.CreateExecutionParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	ctx := context.Background()

	jobID, err := ctl.exec.Start(ctx, execution.Event{
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

	return operations.NewCreateExecutionOK().WithPayload(&operations.CreateExecutionOKBody{
		JobID: &jobID,
	})
}

func (ctl *Execution) Delete(params operations.DeleteExecutionParams) middleware.Responder {
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

func (ctl *Execution) Find(params operations.FindExecutionsParams) middleware.Responder {
	logger := logging.FromRequest(params.HTTPRequest)
	query := dal.Query{
		Pagination: dto.PaginationTo(params.Count, params.Cursor),
		Filtering: dal.Filtering{
			Fields: make([]dal.FilteringField, 0, 2),
		},
	}

	if params.Cause != nil {
		query.Filtering.Fields = append(query.Filtering.Fields, dal.FilteringField{
			Name:  "cause",
			Value: execution.NewCause(*params.Cause),
		})
	}

	if params.Status != nil {
		query.Filtering.Fields = append(query.Filtering.Fields, dal.FilteringField{
			Name:  "status",
			Value: execution.NewStatus(*params.Status),
		})
	}

	ctx := context.Background()
	out, err := ctl.history.Find(ctx, params.ProjectID, query)

	if err != nil {
		logger.Error().
			Timestamp().
			Err(err).
			Str("project_id", params.ProjectID).
			Uint64("count", query.Pagination.Count).
			Str("cursor", query.Pagination.Cursor.String()).
			Msg("failed to find jobs")

		return http.InternalError()
	}

	data := make([]*models.ExecutionOutput, 0, len(out.Data))

	for _, r := range out.Data {
		data = append(data, &models.ExecutionOutput{
			ExecutionCommon: ctl.toExecutionCommonDto(r),
		})
	}

	return operations.NewFindExecutionsOK().WithPayload(&operations.FindExecutionsOKBody{
		Data:         data,
		SearchResult: dto.SearchResultFrom(out.QueryResult),
	})
}

func (ctl *Execution) Get(params operations.GetExecutionParams) middleware.Responder {
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

	return operations.NewGetExecutionOK().WithPayload(&models.ExecutionOutputDetailed{
		ExecutionOutput: models.ExecutionOutput{
			ExecutionCommon: ctl.toExecutionCommonDto(found),
		},
		StartedAt: found.StartedAt.String(),
		EndedAt:   found.EndedAt.String(),
		Params:    dto.ExecutionParamsFrom(found.Params),
	})
}

func (ctl *Execution) toExecutionCommonDto(record history.RecordEntity) models.ExecutionCommon {
	return models.ExecutionCommon{
		ScriptID:  &record.ScriptID,
		ScriptRev: &record.ScriptRev,
		JobID:     &record.JobID,
		Status:    dto.ExecutionStatusFrom(record.Status),
		Cause:     dto.ExecutionCauseFrom(record.CausedBy),
	}
}
