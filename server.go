package home

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Server struct{}

func (s *Server) ListenAndServe() error {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(s.serveHttp)

	return group.Wait()
}
