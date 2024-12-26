package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				log := ExtractLogger(r.Context())

				// TODO: add request id support
				log.Error("Error during request", slog.Any("err", rvr), slog.Any("stack", string(debug.Stack())))

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}

				// TODO: do something on error
			}
		}()

		next.ServeHTTP(w, r)
	})
}
