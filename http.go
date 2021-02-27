package home

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) serveHttp() error {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/{year:[0-9]{4}}/{month:[0-9]{2}}/{slug}", http.HandlerFunc(s.httpPostHandler))
	mux.Get("/{year:[0-9]{4}}/{month:[0-9]{2}}/{slug}/*", http.HandlerFunc(s.httpPostHandler))

	return http.ListenAndServe(":8080", mux)
}

func (s *Server) httpPostHandler(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))
	targetSlug := chi.URLParam(r, "slug")
	targetPath := path.Clean(chi.URLParam(r, "*"))

	// fmt.Println(year, month, targetSlug, targetPath)

	post, err := s.readPost(path.Join("blog", targetSlug))
	if errors.Is(err, fs.ErrNotExist) || post == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		panic(err.Error())
	}

	if post.Date.Year() != year || int(post.Date.Month()) != month {
		basePath := path.Clean(fmt.Sprintf("/%d/%02d/%s", post.Date.Year(), post.Date.Month(), targetSlug))
		if targetPath == "." {
			w.Header().Add("Location", basePath+"/")
		} else {
			w.Header().Add("Location", path.Join(basePath, targetPath))
		}

		w.WriteHeader(http.StatusPermanentRedirect)
		return
	}

	// If we're accessing the post directly, display it.
	if targetPath == "." {
		// All posts need to be rooted at a directory so relative links will
		// work.
		if !strings.HasSuffix(r.URL.Path, "/") {
			w.Header().Add("Location", r.URL.Path+"/")
			w.WriteHeader(http.StatusPermanentRedirect)
			return
		}

		fmt.Fprintf(w, "<h1>%s</h1>\n\n", post.Title)
		fmt.Fprint(w, post.HtmlContent)

		return
	}

	// If we're accessing a path, look up the fs.FS associated with this post
	// and try to find the file.
	targetFS, err := s.lookupPostFS(&post.PostMetadata)
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
