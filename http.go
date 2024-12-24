package home

import (
	"context"
	"net/http"
	"os"

	"github.com/belak/home/internal/middleware"
	"github.com/belak/home/templates"
)

var staticFS = os.DirFS("static")

func (s *Server) serveHttp(ctx context.Context) error {
	mux := http.NewServeMux()

	globalChain := middleware.NewChain(
		middleware.Logger(s.config.Logger),
		middleware.Recoverer(s.config.Logger),
	)

	/*
		requireLogin := internal.NewChain(
			auth.RequireLogin("/login"),
		)
	*/

	mux.HandleFunc("/", s.httpIndexHandler)
	mux.HandleFunc("GET /login", s.httpLoginHandler)
	mux.HandleFunc("POST /login", s.httpLoginHandler)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	return http.ListenAndServe(s.config.BindAddr, globalChain.Then(mux.ServeHTTP))
}

func (s *Server) httpNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.NotFoundHandler(r.URL.Path).Render(r.Context(), w)
}

func (s *Server) httpIndexHandler(w http.ResponseWriter, r *http.Request) {
	// As a special case, if we're in the index handler, and not at the "/"
	// path, we need to call the not found handler.
	if r.URL.Path != "/" {
		s.httpNotFoundHandler(w, r)
		return
	}

	templates.IndexPage().Render(r.Context(), w)
}

func (s *Server) httpLoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.LoginPage().Render(r.Context(), w)
}
