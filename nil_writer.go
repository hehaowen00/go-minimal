package gominimal

import (
	"net/http"
	"strconv"
)

type nilWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func newNilWriter(w http.ResponseWriter) *nilWriter {
	return &nilWriter{
		ResponseWriter: w,
		statusCode:     200,
	}
}

func (w *nilWriter) Write(data []byte) (int, error) {
	n := len(data)
	w.Header().Set("Content-Length", strconv.Itoa(n))
	w.ResponseWriter.Write(nil)
	return n, nil
}
