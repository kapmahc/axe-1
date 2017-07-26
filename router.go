package axe

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter create router
func NewRouter() *Router {
	return &Router{
		routes: make([]route, 0),
	}
}

// Router http router
type Router struct {
	routes []route
}

// GET http get
func (p *Router) GET(path string, handler http.HandlerFunc) {
	p.add(http.MethodGet, path, handler)
}

// POST http post
func (p *Router) POST(path string, handler http.HandlerFunc) {
	p.add(http.MethodGet, path, handler)
}

// PUT http put
func (p *Router) PUT(path string, handler http.HandlerFunc) {
	p.add(http.MethodPut, path, handler)
}

// PATCH http patch
func (p *Router) PATCH(path string, handler http.HandlerFunc) {
	p.add(http.MethodPatch, path, handler)
}

// DELETE http delete
func (p *Router) DELETE(path string, handler http.HandlerFunc) {
	p.add(http.MethodDelete, path, handler)
}

func (p *Router) add(method, path string, handler http.HandlerFunc) {
	p.routes = append(
		p.routes,
		route{
			method:  method,
			path:    path,
			handler: handler,
		},
	)
}

// Group sub-router
func (p *Router) Group(path string, router *Router) {
	for _, rt := range router.routes {
		rt.path = path + rt.path
		p.routes = append(p.routes, rt)
	}
}

// Handler http handle
func (p *Router) Handler() http.Handler {
	rt := mux.NewRouter()
	for _, r := range p.routes {
		rt.HandleFunc(r.path, r.handler).Methods(r.method)
	}
	return rt
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}
