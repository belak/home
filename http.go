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

	mux.Get("/", http.HandlerFunc(s.httpIndexHandler))
	mux.Get("/*", http.HandlerFunc(s.pageHandler))
	mux.Get("/posts", http.HandlerFunc(s.postsHandler))
	mux.Get("/tags", http.HandlerFunc(s.tagsHandler))
	mux.Get("/tags/*", http.HandlerFunc(s.tagsHandler))
	mux.Get("/{year:[0-9]+}/{month:[0-9]+}/{slug}", http.HandlerFunc(s.httpArticleHandler))
	mux.Get("/{year:[0-9]+}/{month:[0-9]+}/{slug}/*", http.HandlerFunc(s.httpArticleHandler))

	staticFS, _ := fs.Sub(s.Content, "static")
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

func (s *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("posts")
}

func (s *Server) tagsHandler(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Clean(chi.URLParam(r, "*"))
	fmt.Println("tags", targetPath)
}

func (s *Server) pageHandler(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Clean(chi.URLParam(r, "*"))
	fmt.Println("page", targetPath)
}

func (s *Server) httpArticleHandler(w http.ResponseWriter, r *http.Request) {
	//year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	//month, _ := strconv.Atoi(chi.URLParam(r, "month"))
	targetSlug := chi.URLParam(r, "slug")
	targetPath := path.Clean(chi.URLParam(r, "*"))

	// fmt.Println(year, month, targetSlug, targetPath)

	article, err := s.readArticle(path.Join("blog", targetSlug))
	if errors.Is(err, fs.ErrNotExist) || article == nil {
		s.httpNotFoundHandler(w, r)
		return
	}

	if err != nil {
		panic(err.Error())
	}

	basePath := path.Clean(fmt.Sprintf("/%d/%02d/%s/", article.Meta.Date.Year(), article.Meta.Date.Month(), targetSlug))

	if !strings.HasPrefix(r.URL.Path, basePath) {
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

		fmt.Fprintf(w, "<h1>%s</h1>\n\n", article.Meta.Title)
		fmt.Fprint(w, article.HtmlContent)

		return
	}

	// If we're accessing a path, look up the fs.FS associated with this post
	// and try to find the file.
	targetFS, err := s.lookupArticleFS(article.Meta)
	if err != nil {
		panic(err.Error())
	}

	f, err := targetFS.Open(targetPath)

	// If the failure was because the file doesn't exist, this counts as not
	// found.
	if errors.Is(err, fs.ErrNotExist) {
		s.httpNotFoundHandler(w, r)
		return
	}

	if err != nil {
		panic(err.Error())
	}

	_, _ = io.Copy(w, f)
}
