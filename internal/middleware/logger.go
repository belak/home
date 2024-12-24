package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				elapsed := time.Since(t1)
				status := ww.Status()
				msg := fmt.Sprintf("Response: %d %s", status, statusLabel(status))

				// TODO: add request id support
				log.With(
					slog.Group("httpRequest",
						slog.String("host", r.Host),
						slog.String("path", r.URL.Path),
					),
					slog.Group("httpResponse",
						slog.Int("status", status),
						slog.Int("bytes", ww.BytesWritten()),
						slog.Duration("elapsed", elapsed),
					),
				).Log(context.Background(), statusLevel(status), msg)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func statusLevel(status int) slog.Level {
	switch {
	case status <= 0:
		return slog.LevelWarn
	case status >= 500:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}
