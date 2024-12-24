package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type contextKey string

const (
	TemplateSetContextKey contextKey = "TemplateSet"
	LoggerContextKey      contextKey = "Logger"
	CurrentUserContextKey contextKey = "CurrentUser"
)

func (c contextKey) String() string {
	return fmt.Sprintf("ContextKey<%s>", string(c))
}

func ExtractLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger); ok {
		return logger
	}

	panic("no logger in context")
}

func ExtractTemplates(ctx context.Context) *TemplateSet {
	if ts, ok := ctx.Value(TemplateSetContextKey).(*TemplateSet); ok {
		return ts
	}
	panic("no template set in context")
}

// contextValueMiddleware is a convenience function which stores a static value in
// the request context using the given contextKey.
func contextValueMiddleware(key contextKey, val interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, key, val)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
