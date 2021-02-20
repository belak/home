package home

import "gopkg.in/gemini.v0"

type wrappedWriter struct {
	bytes  int
	status int
	meta   string

	inner gemini.ResponseWriter
}

func (ww *wrappedWriter) WriteStatus(status int, meta string) {
	if ww.status == 0 {
		ww.status = status
		ww.meta = meta
	}

	ww.inner.WriteStatus(status, meta)
}

func (ww *wrappedWriter) Write(data []byte) (int, error) {
	if ww.status == 0 {
		ww.WriteStatus(gemini.StatusSuccess, "text/gemini")
	}

	n, err := ww.inner.Write(data)
	ww.bytes += n
	return n, err
}
