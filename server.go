package home

import (
	"context"
	"crypto/tls"
	"io/fs"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	TLS     *tls.Config
	Content fs.FS
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListenAndServe() error {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(s.serveHttp)

	return group.Wait()
}
