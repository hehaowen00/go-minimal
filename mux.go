package gominimal

import "net/http"

type customMux struct {
	getHandler     *http.ServeMux
	postHandler    *http.ServeMux
	putHandler     *http.ServeMux
	patchHandler   *http.ServeMux
	deleteHandler  *http.ServeMux
	headHandler    *http.ServeMux
	optionsHandler *http.ServeMux
	connectHandler *http.ServeMux
}

func (mux *customMux) Handle(method string, path string, handler http.Handler) {
	switch method {
	case http.MethodGet:
		mux.getHandler = addRoute(mux.getHandler, path, handler)
	case http.MethodPost:
		mux.postHandler = addRoute(mux.postHandler, path, handler)
	case http.MethodPut:
		mux.putHandler = addRoute(mux.putHandler, path, handler)
	case http.MethodPatch:
		mux.patchHandler = addRoute(mux.patchHandler, path, handler)
	case http.MethodDelete:
		mux.deleteHandler = addRoute(mux.deleteHandler, path, handler)
	case http.MethodHead:
		mux.headHandler = addRoute(mux.headHandler, path, handler)
	case http.MethodOptions:
		mux.optionsHandler = addRoute(mux.optionsHandler, path, handler)
	case http.MethodConnect:
		mux.connectHandler = addRoute(mux.connectHandler, path, handler)
	default:
		panic("invalid route method")
	}
}

func (mux *customMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		checkHandler(mux.getHandler, w, r)
	case http.MethodPost:
		checkHandler(mux.postHandler, w, r)
	case http.MethodPut:
		checkHandler(mux.putHandler, w, r)
	case http.MethodPatch:
		checkHandler(mux.patchHandler, w, r)
	case http.MethodDelete:
		checkHandler(mux.deleteHandler, w, r)
	case http.MethodHead:
		w = newNilWriter(w)
		checkHandler(mux.getHandler, w, r)
	case http.MethodConnect:
		checkHandler(mux.connectHandler, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
}

func addRoute(mux *http.ServeMux, path string, handler http.Handler) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux()
	}
	mux.Handle(path, handler)
	return mux
}

func checkHandler(mux *http.ServeMux, w http.ResponseWriter, r *http.Request) {
	if mux == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	mux.ServeHTTP(w, r)
}
