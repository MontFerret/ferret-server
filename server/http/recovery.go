package http

import (
	"github.com/MontFerret/ferret-server/server/logging"
	"github.com/pkg/errors"
	"net/http"
	"runtime/debug"
)

type Recovery struct{}

func NewRecovery() *Recovery {
	return &Recovery{}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer func() {
		logger := logging.FromRequest(req)

		if r := recover(); r != nil {
			var err error

			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(""))

			logger.Error().
				Timestamp().
				Err(err).
				Str("stack", string(debug.Stack())).
				Msgf("PANIC")
		}
	}()

	next(rw, req)
}
