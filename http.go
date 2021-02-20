package home

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) serveHttp() error {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/{year:[0-9]{4}}/{month:[0-9]{2}}/{slug}/*", http.HandlerFunc(s.httpPostHandler))

	return http.ListenAndServe(":8080", mux)
}

func (s *Server) httpPostHandler(w http.ResponseWriter, r *http.Request) {
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")
	slug := chi.URLParam(r, "slug")
	targetPath := path.Clean(chi.URLParam(r, "*"))

	fmt.Println(year, month, slug, targetPath)

	post := s.lookupCachedPost(slug)

	// Nil post means not found, so we return early
	if post == nil {
		return
	}

	if targetPath == "." {
		// All posts need to be rooted at a directory so relative links will
		// work.
		if !strings.HasSuffix(r.URL.Path, "/") {
			w.Header().Add("Location", r.URL.Path+"/")
			w.WriteHeader(http.StatusPermanentRedirect)
			return
		}

		fmt.Fprintf(w, "# %s\n\n", post.Meta.DisplayTitle())

		_, err := io.WriteString(w, post.Body)
		if err != nil {
			panic(err.Error())
		}

		return
	}

	targetFS, err := s.lookupPostFS(post.Meta)
	if err != nil {
		panic(err.Error())
	}

	f, err := targetFS.Open(targetPath)

	// If the failure was because the file doesn't exist, this counts as not
	// found.
	if errors.Is(err, fs.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		panic(err.Error())
	}

	_, _ = io.Copy(w, f)
}
