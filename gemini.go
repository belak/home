package home

import (
	"context"
	"os"

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
	}

	w.WriteStatus(gemini.StatusInput, "something: "+params[0])
}
