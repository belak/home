package home

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type ServerConfig struct {
	Logger   *slog.Logger
	BindAddr string
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		config: config,
	}
}

type Server struct {
	config ServerConfig
}

func (s *Server) ListenAndServe() error {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() {
		return s.serveHttp(ctx)
	})

	return group.Wait()
}
