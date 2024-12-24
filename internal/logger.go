package internal

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
)

type LoggerConfig struct {
	Level slog.Level `envconfig:"level"`
}

func NewLogger() (*slog.Logger, error) {
	var logger *slog.Logger

	config := LoggerConfig{
		Level: slog.LevelDebug,
	}

	err := envconfig.Process("LOG", &config)

	if Env() == EnvDev {
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: config.Level}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.Level}))
	}

	return logger, err
}

func InjectLoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return contextValueMiddleware(LoggerContextKey, logger)
}
