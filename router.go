package axe

import (
	"html/template"
	"net/http"
	"sort"
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
		layouts:  make(map[string]LayoutFunc),
	}
}

// Router http router
type Router struct {
	path     string
	routes   []route
	handlers []HandlerFunc
	funcs    template.FuncMap
	statics  map[string]string
	layouts  map[string]LayoutFunc
}

// Use use handlers
func (p *Router) Use(handlers ...HandlerFunc) {
	p.handlers = append(p.handlers, handlers...)
}

// FuncMap set html template funcmap
func (p *Router) FuncMap(n string, f interface{}) {
	p.funcs[n] = f
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
	p.add(http.MethodPost, path, handlers...)
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

// Resources resources
func (p *Router) Resources(path string, index, create, show, update, destroy []HandlerFunc) {
	if index != nil {
		p.GET(path, index...)
	}
	if create != nil {
		p.POST(path, create...)
	}

	path += "/{id}"
	if show != nil {
		p.GET(path, show...)
	}
	if update != nil {
		p.POST(path, update...)
	}
	if destroy != nil {
		p.POST(path, destroy...)
	}
}

// Group sub-router
func (p *Router) Group(path string, router *Router) {
	for _, rt := range router.routes {
		p.add(
			rt.method,
			path+rt.path,
			rt.handlers...,
		)
	}
	for k, v := range router.statics {
		p.statics[k] = v
	}
	for k, v := range router.layouts {
		p.layouts[k] = v
	}
	for k, v := range router.funcs {
		p.funcs[k] = v
	}
}

// Walk walk routes
func (p *Router) Walk(f func(method, path string, handlers ...HandlerFunc) error) error {
	sort.Sort(routes(p.routes))
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
		log.Debugf("GET %s => %s", k, v)
		rut.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v)))).Methods(http.MethodGet)
	}
	for _, r := range p.routes {
		log.Debugf("%s %s [%d]", r.method, r.path, len(r.handlers))
		rut.HandleFunc(r.path, func(wrt http.ResponseWriter, req *http.Request) {
			now := time.Now()
			log.Infof("%s %s %s", req.Proto, req.Method, req.URL)
			log.Debug(req.Header)
			ctx := Context{
				Request:  req,
				Writer:   wrt,
				Payload:  make(H),
				Params:   mux.Vars(req),
				index:    0,
				handlers: r.handlers,
				decoder:  dec,
				validate: vat,
				render:   rdr,
			}
			ctx.Next()
			log.Infof("%d %s", ctx.code, time.Now().Sub(now))
		}).Methods(r.method)
	}
	return rut
}
