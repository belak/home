package home

import (
	"context"
	"crypto/tls"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	TLS *tls.Config
}

func (s *Server) ListenAndServe() error {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(s.serveHttp)
	group.Go(s.serveGemini)

	return group.Wait()
}
