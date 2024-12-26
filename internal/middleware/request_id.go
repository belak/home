package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

var requestIDContextKey ContextKey = "RequestID"

func InjectRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := ExtractLogger(r.Context())

		requestID, err := uuid.NewV7()
		if err != nil {
			log.Warn("Failed to generate Request ID", slog.Any("err", err))
		} else {
			ctx := context.WithValue(r.Context(), requestIDContextKey, requestID.String())
			r = r.WithContext(ctx)

			log.With(slog.String("requestID", requestID.String()))
			w.Header().Set("Request-ID", requestID.String())
		}

		next.ServeHTTP(w, r)
	})
}

func ExtractRequestID(ctx context.Context) string {
	ret, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}

	return ret
}
