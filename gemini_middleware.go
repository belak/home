package home

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"gopkg.in/gemini.v0"
)

func geminiRecoverer(next gemini.Handler) gemini.Handler {
	return gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(os.Stderr, "panic: %v\n\n%s\n", err, debug.Stack())
				w.WriteStatus(gemini.StatusCGIError, "internal error")
			}
		}()

		next.ServeGemini(ctx, w, r)
	})
}

func geminiLogger(next gemini.Handler) gemini.Handler {
	return gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		buf := &bytes.Buffer{}
		start := time.Now()

		ww := &wrappedWriter{inner: w}

		defer func() {
			if ww.status == 0 {
				ww.WriteStatus(gemini.StatusNotFound, "not found")
			}

			elapsed := time.Since(start)

			cW(buf, nCyan, "\"%s %s GEMINI/0.14.3\"", r.ServerName, r.URL)
			buf.WriteString(" from ")
			buf.WriteString(r.RemoteAddr)
			buf.WriteString(" - ")

			status := ww.status

			switch (status / 10) * 10 {
			case 10: // Input
				cW(buf, bBlue, "%02d", status)
			case 20: // Success
				cW(buf, bGreen, "%02d", status)
			case 30: // Redirect
				cW(buf, bCyan, "%02d", status)
			case 40: // Temporary Failure
				cW(buf, bYellow, "%02d", status)
			case 50: // Permanent Failure
				cW(buf, bRed, "%02d", status)
			default: // Certificate or Other
				cW(buf, bWhite, "%02d", status)
			}

			fmt.Fprintf(buf, " %q ", ww.meta)

			cW(buf, bBlue, "%dB", ww.bytes)

			buf.WriteString(" in ")
			if elapsed < 500*time.Millisecond {
				cW(buf, nGreen, "%s", elapsed)
			} else if elapsed < 5*time.Second {
				cW(buf, nYellow, "%s", elapsed)
			} else {
				cW(buf, nRed, "%s", elapsed)
			}

			log.Print(buf.String())
		}()

		next.ServeGemini(ctx, ww, r)
	})
}
