package dto

import (
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/server/http/api/models"
)

func ExecutionCauseFrom(cause execution.Cause) models.ExecutionCause {
	switch cause {
	case execution.CauseManual:
		return models.ExecutionCauseManual
	case execution.CauseSchedule:
		return models.ExecutionCauseSchedule
	case execution.CauseHook:
		return models.ExecutionCauseHook
	default:
		return models.ExecutionCauseUnknown
	}
}

func ExecutionCauseTo(cause models.ExecutionCause) execution.Cause {
	switch cause {
	case models.ExecutionCauseManual:
		return execution.CauseManual
	case models.ExecutionCauseSchedule:
		return execution.CauseSchedule
	case models.ExecutionCauseHook:
		return execution.CauseHook
	default:
		return execution.CauseUnknown
	}
}
