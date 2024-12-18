package home

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/belak/home/templates"
)

var staticFS = os.DirFS("static")

func (s *Server) serveHttp(ctx context.Context) error {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", http.HandlerFunc(s.httpIndexHandler))

	mux.Mount("/static", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	mux.NotFound(s.httpNotFoundHandler)

	return http.ListenAndServe(s.config.BindAddr, mux)
}

func (s *Server) httpNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Println("not found")
	templates.PageErrNotFound(r.URL.Path).Render(r.Context(), w)
}

func (s *Server) httpIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index")
}
