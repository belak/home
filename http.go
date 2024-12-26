package home

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/belak/home/internal"
	"github.com/belak/home/internal/middleware"
	"github.com/belak/home/models"
	"github.com/belak/home/templates"
)

func mustSubFS(rawFS fs.FS, dir string) fs.FS {
	ret, err := fs.Sub(rawFS, dir)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

//go:embed static
var rawStaticFS embed.FS
var staticFS fs.FS = mustSubFS(rawStaticFS, "static")

func (s *Server) serveHttp(ctx context.Context) error {
	mux := http.NewServeMux()

	globalChain := middleware.NewChain(
		middleware.InjectLogger(s.config.Logger),
		middleware.InjectRequestID,
		middleware.RequestLogger,
		middleware.Recoverer,
	)

	/*
		requireLogin := internal.NewChain(
			auth.RequireLogin("/login"),
		)
	*/

	mux.HandleFunc("/", s.httpIndexHandler)
	mux.HandleFunc("GET  /login", s.httpLoginHandler)
	mux.HandleFunc("POST /login", s.httpLoginHandler)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	return http.ListenAndServe(s.config.BindAddr, globalChain.Then(mux.ServeHTTP))
}

func (s *Server) httpNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.NotFoundHandler(&models.SiteState{}, r.URL.Path).Render(r.Context(), w)
}

func (s *Server) httpIndexHandler(w http.ResponseWriter, r *http.Request) {
	// As a special case, if we're in the index handler, and not at the "/"
	// path, we need to call the not found handler.
	if r.URL.Path != "/" {
		s.httpNotFoundHandler(w, r)
		return
	}

	templates.IndexPage(&models.SiteState{}).Render(r.Context(), w)
}

func (s *Server) httpLoginHandler(w http.ResponseWriter, r *http.Request) {
	var form templates.LoginForm
	data, ok := internal.Bind(r, &form)
	fmt.Println(data, ok)
	if !ok {
		// TODO: bad request?
		return
	}
	templates.LoginPage(&models.SiteState{}, &form).Render(r.Context(), w)
}
