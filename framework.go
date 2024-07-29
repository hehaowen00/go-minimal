package gominimal

import (
	"net/http"
	"strings"
)

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

type Router struct {
	mux        *customMux
	middleware []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{
		mux: &customMux{},
	}
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append([]MiddlewareFunc{middleware}, r.middleware...)
}

func (r *Router) Serve(addr string) error {
	mux := r.mux
	r.mux = nil
	return http.ListenAndServe(addr, mux)
}

func (r *Router) ServeTLS(addr string, cert, key string) error {
	mux := r.mux
	r.mux = nil
	return http.ListenAndServeTLS(addr, cert, key, mux)
}

func (r *Router) GET(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	for i := range middleware {
		handler = middleware[len(middleware)-1-i](handler)
	}

	for _, h := range r.middleware {
		handler = h(handler)
	}

	r.route(http.MethodGet, path, handler)
}

func (r *Router) route(method string, path string, handler http.HandlerFunc) {
	path = strings.TrimSpace(path)
	r.mux.Handle(method, path, handler)
}

type customMux struct {
	getHandler    *http.ServeMux
	postHandler   *http.ServeMux
	putHandler    *http.ServeMux
	deleteHandler *http.ServeMux
}

func (mux *customMux) Handle(method string, path string, handler http.Handler) {
	switch method {
	case http.MethodGet:
		if mux.getHandler == nil {
			mux.getHandler = http.NewServeMux()
		}
		mux.getHandler.Handle(path, handler)
	case http.MethodPost:
		if mux.postHandler == nil {
			mux.postHandler = http.NewServeMux()
		}
		mux.postHandler.Handle(path, handler)
	}
}

func (mux *customMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if mux.getHandler == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
		}
		mux.getHandler.ServeHTTP(w, r)
	case http.MethodPost:
		if mux.postHandler == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
		}
		mux.postHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
}
