package gominimal

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"
)

type CorsOptions struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

func CorsMiddleware(opts CorsOptions) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if len(opts.AllowedOrigins) == 0 {
				origin := r.Header.Get("Origin")
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			} else {
				w.Header().Set(
					"Access-Control-Allow-Origin",
					strings.Join(opts.AllowedOrigins, ","),
				)
			}

			if opts.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			next(w, r)
		}
	}
}

func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Accept-Encoding")

		if !strings.Contains(val, "gzip") {
			next(w, r)
			return
		}

		writer := gzipWriter{
			ResponseWriter: w,
		}

		writer.writer = &bufWriter{
			parent: &writer,
		}

		next(&writer, r)
		writer.Close()
	}
}

type writerStrategy interface {
	Close() error
	Write([]byte) (int, error)
}

type bufWriter struct {
	parent *gzipWriter
	buf    bytes.Buffer
}

func (w *bufWriter) Write(data []byte) (int, error) {
	n, err := w.buf.Write(data)

	if w.buf.Len() >= 250 {
		w.parent.useGzip()
		_, err = w.parent.writer.Write(w.buf.Bytes())
		w.buf.Reset()
	}

	return n, err
}

func (w *bufWriter) Close() error {
	w.parent.ResponseWriter.Write(w.buf.Bytes())
	return nil
}

type gzipWriter struct {
	http.ResponseWriter
	writer writerStrategy
}

func (w *gzipWriter) useGzip() {
	w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	w.writer = gzip.NewWriter(w.ResponseWriter)
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *gzipWriter) Close() {
	w.writer.Close()
}
