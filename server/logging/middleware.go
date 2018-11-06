package logging

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger    *zerolog.Logger
	headerKey string
}

func NewMiddleware(logger *zerolog.Logger, headerKey string) *LoggerMiddleware {
	return &LoggerMiddleware{logger, headerKey}
}

func (l *LoggerMiddleware) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	id := req.Header.Get(l.headerKey)

	start := time.Now()
	path := req.URL.Path
	query := req.URL.RawQuery

	reqLogger := l.logger.With().Str("request_id", id).Logger()
	nextCtx := reqLogger.WithContext(req.Context())
	next.ServeHTTP(resp, req.WithContext(nextCtx))

	end := time.Now()
	latency := end.Sub(start)

	reqLogger.Info().
		Timestamp().
		Str("method", req.Method).
		Str("path", path).
		Str("query", query).
		Str("ip", req.UserAgent()).
		Str("user-agent", req.UserAgent()).
		Str("time", end.Format(time.RFC3339)).
		Dur("latency", latency).
		Msg("Incoming request")
}
