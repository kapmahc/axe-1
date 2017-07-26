package axe

import (
	"net/http"

	"github.com/unrolled/render"
)

// Wrapper wrapper
type Wrapper struct {
	R render.Render
}

func (p *Wrapper) error(w http.ResponseWriter, e error) {
	s := http.StatusInternalServerError
	if er, ok := e.(*HTTPError); ok {
		s = er.Status
	}
	http.Error(w, e.Error(), s)
}

// HTML rener html
func (p *Wrapper) HTML(l string, n string, f func(*Context) (H, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, e := f(newContext(w, r))
		if e == nil {
			v["error"] = e.Error()
		}
		p.R.HTML(w, http.StatusOK, n, v, render.HTMLOptions{Layout: l})
	}
}

// JSON rener json
func (p *Wrapper) JSON(f func(*Context) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, e := f(newContext(w, r))
		if e == nil {
			p.R.JSON(w, http.StatusOK, v)
			return
		}
		p.error(w, e)
	}
}

// XML rener xml
func (p *Wrapper) XML(f func(*Context) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, e := f(newContext(w, r))
		if e == nil {
			p.R.XML(w, http.StatusOK, v)
			return
		}
		p.error(w, e)
	}
}
