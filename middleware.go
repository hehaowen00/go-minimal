package gominimal

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CorsOptions struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

func CorsMiddleware(opts CorsOptions, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(opts.AllowedOrigins) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
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

func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Accept-Encoding")

		if !strings.Contains(val, "gzip") {
			next(w, r)
			return
		}

		writer := tempWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(&writer, r)

		w.WriteHeader(writer.statusCode)

		if writer.buf.Len() < 250 {
			io.Copy(w, &writer.buf)
			return
		}

		w.Header().Add("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()
		io.Copy(gz, &writer.buf)
	}
}

type tempWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func (w *tempWriter) WriteHeader(status int) {
	w.statusCode = status
}

func (w *tempWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}
