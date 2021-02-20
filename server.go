package home

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	TLS     *tls.Config
	Content fs.FS

	cachedPostsLock sync.RWMutex
	cachedPosts     map[string]*cachedPost
}

func NewServer() *Server {
	return &Server{
		cachedPosts: make(map[string]*cachedPost),
	}
}

func (s *Server) ListenAndServe() error {
	group, _ := errgroup.WithContext(context.Background())

	err := s.updateCachedPosts()
	if err != nil {
		fmt.Println(err)
	}

	group.Go(s.serveHttp)
	group.Go(s.serveGemini)

	return group.Wait()
}
