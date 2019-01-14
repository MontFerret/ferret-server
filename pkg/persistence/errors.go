package persistence

import "github.com/pkg/errors"

var (
	ErrNoPersistence = errors.New("local persistence not available")
)
