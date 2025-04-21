package gominimal

import (
	"net/http"
	"net/url"
)

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

type router struct {
	mux        *customMux
	middleware []MiddlewareFunc
	basePath   string
}

type IRouter interface {
	Get(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	Post(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	Patch(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	Put(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	Delete(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)

	// Head(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	// Options(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)
	// Connect(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc)

	Group(path string, subgroup ...bool) IRouter
}

func NewRouter() *router {
	return &router{
		mux: &customMux{
			optionsMap: make(map[string][]string),
		},
	}
}

func (r *router) Handler() http.Handler {
	mux := r.mux
	mux.buildOptions()
	r.mux = nil
	return mux
}

func (r *router) Use(middleware MiddlewareFunc) {
	r.middleware = append([]MiddlewareFunc{middleware}, r.middleware...)
}

func (r *router) Get(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodGet, r.basePath, path, handler)
}

func (r *router) Post(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPost, r.basePath, path, handler)
}

func (r *router) Put(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPut, r.basePath, path, handler)
}

func (r *router) Patch(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodPatch, r.basePath, path, handler)
}

func (r *router) Delete(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	handler = applyMiddleware(handler, middleware, r.middleware)
	r.route(http.MethodDelete, r.basePath, path, handler)
}

// func (r *router) Head(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
// 	handler = applyMiddleware(handler, middleware, r.middleware)
// 	r.route(http.MethodGet, r.basePath, path, handler)
// }

// func (r *router) Options(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
// 	handler = applyMiddleware(handler, middleware, r.middleware)
// 	r.route(http.MethodGet, r.basePath, path, handler)
// }

// func (r *router) Connect(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
// 	handler = applyMiddleware(handler, middleware, r.middleware)
// 	r.route(http.MethodGet, r.basePath, path, handler)
// }

func (r *router) Group(path string, subgroup ...bool) IRouter {
	var middlewares []MiddlewareFunc

	if len(subgroup) == 1 && subgroup[0] {
		path, _ = url.JoinPath(r.basePath, path)
		middlewares = r.middleware
	}

	newRouter := *r
	newRouter.basePath = path
	newRouter.middleware = middlewares

	return &newRouter
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

func (r *router) route(method string, basePath string, path string, handler http.HandlerFunc) {
	path, err := url.JoinPath(basePath, path)
	if err != nil {
		panic(err)
	}

	r.mux.handle(method, path, handler)
}
