package router

import "net/http"

type Group struct {
	prefix          string
	router          *ServeMux
	preMiddlewares  []Middleware
	postMiddlewares []Middleware
}

func (mux *ServeMux) GROUP(prefix string, preMiddlewares []Middleware, postMiddlewares []Middleware) *Group {
	return &Group{
		prefix:          prefix,
		router:          mux,
		preMiddlewares:  preMiddlewares,
		postMiddlewares: postMiddlewares,
	}
}

func (g *Group) addRoute(method, path string, handler Handler, middlewares []Middleware) {
	fullPath := g.prefix + path
	allMiddlewares := append(g.preMiddlewares, middlewares...)
	allMiddlewares = append(allMiddlewares, g.postMiddlewares...)
	finalHandler := applyMiddlewares(handler, allMiddlewares...)
	g.router.ServeMux.Handle(method+" "+fullPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ninaRequest := &NinaRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			UserAgent:     r.UserAgent(),
		}
		finalHandler(w, ninaRequest)
	}))
}

func (g *Group) GET(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodGet, path, handler, middlewares)
}

func (g *Group) POST(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodPost, path, handler, middlewares)
}

func (g *Group) PUT(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodPut, path, handler, middlewares)
}

func (g *Group) DELETE(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodDelete, path, handler, middlewares)
}

func (g *Group) HEAD(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodHead, path, handler, middlewares)
}

func (g *Group) PATCH(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodPatch, path, handler, middlewares)
}

func (g *Group) OPTIONS(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodOptions, path, handler, middlewares)
}

func (g *Group) CONNECT(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodConnect, path, handler, middlewares)
}

func (g *Group) TRACE(path string, handler Handler, middlewares []Middleware) {
	g.addRoute(http.MethodTrace, path, handler, middlewares)
}
