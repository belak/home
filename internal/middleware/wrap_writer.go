package middleware

// The original work was derived from Chi's middleware which was in turn derived
// from Goji's middleware, sources:
//
// https://github.com/go-chi/chi/blob/master/middleware/wrap_writer.go
// https://github.com/zenazn/goji/tree/master/web/middleware
//
// However, this version has been simplified - it does less, allows for less
// introspection and tries to accomplish its job as simply as possible, falling
// back to

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// NewWrapResponseWriter wraps an http.ResponseWriter, returning a proxy that
// allows you to introspect various information about the response.
func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int) WrapResponseWriter {
	_, fl := w.(http.Flusher)

	bw := basicWriter{ResponseWriter: w}

	if protoMajor == 2 {
		// The http/2 ResponseWriter from net/http supports http.Flusher and
		// http.Pusher, but specifically doesn't support http.Hijacker.
		_, ps := w.(http.Pusher)
		if fl && ps {
			return &http2FancyWriter{bw}
		}
	} else {
		// The standard http ResponseWriter from net/http supports http.Flusher,
		// http.Hijacker, and io.ReaderFrom (though the io.ReaderFrom support is
		// not specifically documented).
		//
		// Because the io.ReaderFrom interface is not documented, we want to be
		// careful to fallback to an otherwise fully featured wrapper if it's
		// not supported.
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		if fl && hj {
			if rf {
				return &httpFancyReaderFromWriter{httpFancyWriter{bw}}
			} else {
				return &httpFancyWriter{bw}
			}
		}
	}

	return &bw
}

// WrapResponseWriter is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type WrapResponseWriter interface {
	http.ResponseWriter

	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int

	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int

	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
}

func (b *basicWriter) WriteHeader(code int) {
	if code >= 100 && code <= 199 && code != http.StatusSwitchingProtocols {
		b.ResponseWriter.WriteHeader(code)
	} else if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basicWriter) Write(buf []byte) (n int, err error) {
	b.maybeWriteHeader()
	n, err = b.ResponseWriter.Write(buf)
	b.bytes += n
	return n, err
}

func (b *basicWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

func (b *basicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

// httpFancyWriter is a HTTP writer that additionally satisfies
// http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type httpFancyWriter struct {
	basicWriter
}

func (f *httpFancyWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (f *httpFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (f *httpFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	rf := f.basicWriter.ResponseWriter.(io.ReaderFrom)
	f.basicWriter.maybeWriteHeader()
	n, err := rf.ReadFrom(r)
	f.basicWriter.bytes += int(n)
	return n, err
}

var _ http.Flusher = &httpFancyWriter{}
var _ http.Hijacker = &httpFancyWriter{}

type httpFancyReaderFromWriter struct {
	httpFancyWriter
}

var _ io.ReaderFrom = &httpFancyReaderFromWriter{}

// http2FancyWriter is a HTTP2 writer that additionally satisfies http.Flusher,
// and http.Pusher. It exists for the common case of wrapping the
// http.ResponseWriter that package http gives you for http/2 connections, in
// order to make the proxied object support the full method set of the proxied
// object. Note that this specifically ignores the deprecated http.CloseNotifier interface.
type http2FancyWriter struct {
	basicWriter
}

func (f *http2FancyWriter) Push(target string, opts *http.PushOptions) error {
	return f.basicWriter.ResponseWriter.(http.Pusher).Push(target, opts)
}

func (f *http2FancyWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.Flusher = &http2FancyWriter{}
var _ http.Pusher = &http2FancyWriter{}
