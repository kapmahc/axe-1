package axe

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewRouter create router
func NewRouter() *Router {
	return &Router{
		routes:   make([]route, 0),
		handlers: make([]HandlerFunc, 0),
		funcs:    make(template.FuncMap),
		statics:  make(map[string]string),
	}
}

// Router http router
type Router struct {
	path     string
	routes   []route
	handlers []HandlerFunc
	funcs    template.FuncMap
	statics  map[string]string
}

// Use use handlers
func (p *Router) Use(handlers ...HandlerFunc) {
	p.handlers = append(p.handlers, handlers...)
}

// Static static files
func (p *Router) Static(path, dir string) {
	p.statics[path] = dir
}

// GET http get
func (p *Router) GET(path string, handlers ...HandlerFunc) {
	p.add(http.MethodGet, path, handlers...)
}

// POST http post
func (p *Router) POST(path string, handlers ...HandlerFunc) {
	p.add(http.MethodGet, path, handlers...)
}

// PUT http put
func (p *Router) PUT(path string, handlers ...HandlerFunc) {
	p.add(http.MethodPut, path, handlers...)
}

// PATCH http patch
func (p *Router) PATCH(path string, handlers ...HandlerFunc) {
	p.add(http.MethodPatch, path, handlers...)
}

// DELETE http delete
func (p *Router) DELETE(path string, handlers ...HandlerFunc) {
	p.add(http.MethodDelete, path, handlers...)
}

func (p *Router) add(method, path string, handlers ...HandlerFunc) {
	p.routes = append(
		p.routes,
		route{
			method:   method,
			path:     path,
			handlers: append(p.handlers, handlers...),
		},
	)
}

// Group sub-router
func (p *Router) Group(path string, router *Router) {
	for _, rt := range router.routes {
		rt.path = path + rt.path
		rt.handlers = append(p.handlers, rt.handlers...)
		p.routes = append(p.routes, rt)
	}
	for k, v := range router.statics {
		p.statics[k] = v
	}
}

// Walk walk routes
func (p *Router) Walk(f func(method, path string, handlers ...HandlerFunc) error) error {
	for _, r := range p.routes {
		if e := f(r.method, r.path, r.handlers...); e != nil {
			return e
		}
	}
	return nil
}

// Handler http handle
func (p *Router) Handler(views string, debug bool) http.Handler {
	rdr := render.New(render.Options{
		Directory:     views,
		Extensions:    []string{".html"},
		Funcs:         []template.FuncMap{},
		IndentJSON:    debug,
		IndentXML:     debug,
		IsDevelopment: debug,
	})

	rut := mux.NewRouter()
	dec := form.NewDecoder()
	vat := validator.New()
	for k, v := range p.statics {
		rut.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v)))).Methods(http.MethodGet)
	}
	for _, r := range p.routes {
		rut.HandleFunc(r.path, func(wrt http.ResponseWriter, req *http.Request) {
			now := time.Now()
			log.Info(req.Proto, req.Method, req.URL)
			ctx := Context{
				Request:  req,
				Writer:   wrt,
				Payload:  make(H),
				Params:   mux.Vars(req),
				index:    -1,
				handlers: r.handlers,
				decoder:  dec,
				validate: vat,
				render:   rdr,
			}
			ctx.Next()
			log.Infof("done, %d %s", ctx.code, time.Now().Sub(now))
		}).Methods(r.method)
	}
	return rut
}

type route struct {
	method   string
	path     string
	handlers []HandlerFunc
}
