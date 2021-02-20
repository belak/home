package home

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"gopkg.in/gemini.v0"
)

func (s *Server) serveGemini() error {
	mux := gemini.NewServeMux()

	mux.Register("/:rest", gemini.FS(os.DirFS("content")))
	mux.Register("/posts/:rest", gemini.HandlerFunc(s.geminiPostHandler))

	gem := gemini.Server{
		Handler: geminiLogger(geminiRecoverer(mux)),
		TLS:     s.TLS.Clone(),
		Log:     &gemini.NopServerLogger{},
	}

	return gem.ListenAndServe()
}

func (s *Server) geminiPostHandler(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
	params := gemini.CtxParams(ctx)
	if len(params) != 1 {
		w.WriteStatus(gemini.StatusCGIError, "internal error")
		return
	}

	targetPath := path.Clean(params[0])

	// If we're at the post root, list all posts.
	if targetPath == "." {
		posts := s.lookupPosts(true)

		fmt.Fprintf(w, "# Posts\n\n")

		for _, post := range posts {
			formattedDate := post.Date.Format("2006-01-02")

			fmt.Fprintf(w,
				"=> %s/ %s - %s",
				post.Slug,
				formattedDate,
				post.DisplayTitle())

			if post.Draft {
				fmt.Fprintf(w, " [DRAFT]")
			}

			fmt.Fprint(w, "\n")
		}
		return
	}

	parts := strings.SplitN(targetPath, "/", 2)
	post := s.lookupCachedPost(parts[0])

	// Nil post means not found, so we return early
	if post == nil {
		return
	}

	if len(parts) == 1 {
		// All posts need to be rooted at a directory so relative links will
		// work.
		if !strings.HasSuffix(r.URL.Path, "/") {
			w.WriteStatus(gemini.StatusPermanentRedirect, r.URL.Path+"/")
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

	f, err := targetFS.Open(parts[1])

	// If the failure was because the file doesn't exist, return early so the
	// not found handler kicks in.
	if errors.Is(err, fs.ErrNotExist) {
		return
	}

	if err != nil {
		panic(err.Error())
	}

	_, _ = io.Copy(w, f)
}
