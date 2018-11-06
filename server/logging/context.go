package logging

import (
	"context"
	"github.com/rs/zerolog"
	"net/http"
)

func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func FromRequest(req *http.Request) *zerolog.Logger {
	return FromContext(req.Context())
}
