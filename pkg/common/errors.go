package common

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrMissedArgument        = errors.New("missed argument")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrInvalidArgumentNumber = errors.New("invalid argument number")
	ErrInvalidType           = errors.New("invalid type")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrNotFound              = errors.New("not found")
	ErrNotUnique             = errors.New("not unique")
	ErrTerminated            = errors.New("operation is terminated")
	ErrUnexpected            = errors.New("unexpected error")
	ErrTimeout               = errors.New("operation timed out")
	ErrNotImplemented        = errors.New("not implemented")
)

type AggregatedError struct {
	errors []error
}

func NewAggregatedError(errs []error) AggregatedError {
	return AggregatedError{errs}
}

func (e AggregatedError) Errors() []error {
	return e.errors
}

func (e AggregatedError) Error() string {
	return Errors(e.errors...).Error()
}

func Error(err error, msg string) error {
	return errors.Errorf("%s: %s", err.Error(), msg)
}

func Errorf(err error, format string, args ...interface{}) error {
	return errors.Errorf("%s: %s", err.Error(), fmt.Sprintf(format, args...))
}

func Errors(err ...error) error {
	message := ""

	for _, e := range err {
		message += ": " + e.Error()
	}

	return errors.New(message)
}
