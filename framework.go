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
		mux: &customMux{
			optionsMap: make(map[string][]string),
		},
	}
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append([]MiddlewareFunc{middleware}, r.middleware...)
}

func (r *Router) Serve(addr string) error {
	mux := r.mux
	mux.buildOptions()

	r.mux = nil
	return http.ListenAndServe(addr, mux)
}

func (r *Router) ServeTLS(addr string, cert, key string) error {
	mux := r.mux
	mux.buildOptions()

	r.mux = nil
	return http.ListenAndServeTLS(addr, cert, key, mux)
}

func (r *Router) GET(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodDelete, path, handler)
}

func (r *Router) HEAD(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodHead, path, handler)
}

func (r *Router) OPTIONS(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodOptions, path, handler)
}

func (r *Router) CONNECT(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodConnect, path, handler)
}

func applyMiddleware(
	handler http.HandlerFunc,
	m1 []MiddlewareFunc,
	m2 []MiddlewareFunc,
) http.HandlerFunc {
	for i := range m1 {
		handler = m1[len(m1)-1-i](handler)
	}

	for _, h := range m2 {
		handler = h(handler)
	}

	return handler
}

func (r *Router) route(method string, path string, handler http.HandlerFunc) {
	path = strings.TrimSpace(path)
	r.mux.Handle(method, path, handler)
}
