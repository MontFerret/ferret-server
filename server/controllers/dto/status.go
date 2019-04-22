package dto

import (
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/server/http/api/models"
)

func ExecutionStatusFrom(status execution.Status) models.ExecutionStatus {
	switch status {
	case execution.StatusQueued:
		return models.ExecutionStatusQueued
	case execution.StatusRunning:
		return models.ExecutionStatusRunning
	case execution.StatusCompleted:
		return models.ExecutionStatusCompleted
	case execution.StatusCancelled:
		return models.ExecutionStatusCancelled
	case execution.StatusErrored:
		return models.ExecutionStatusErrored
	default:
		return models.ExecutionStatusUnknown
	}
}

func ExecutionStatusTo(status models.ExecutionStatus) execution.Status {
	switch status {
	case models.ExecutionStatusQueued:
		return execution.StatusQueued
	case models.ExecutionStatusRunning:
		return execution.StatusRunning
	case models.ExecutionStatusCompleted:
		return execution.StatusCompleted
	case models.ExecutionStatusCancelled:
		return execution.StatusCancelled
	case models.ExecutionStatusErrored:
		return execution.StatusErrored
	default:
		return execution.StatusUnknown
	}
}
