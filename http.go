package home

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var staticFS = os.DirFS("static")
var templatesFS = os.DirFS("templates")

func (s *Server) serveHttp() error {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", http.HandlerFunc(s.httpIndexHandler))

	mux.Mount("/static", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	mux.NotFound(s.httpNotFoundHandler)

	return http.ListenAndServe(":8080", mux)
}

func (s *Server) httpNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.httpExecuteTemplate(w, r, "404.html", "", &NotFoundContext{Path: r.URL.Path})
}

func (s *Server) httpIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index")
}
