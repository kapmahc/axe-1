package axe

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/go-playground/form"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// HandlerFunc http handler
type HandlerFunc func(*Context)

// H hash data
type H map[string]interface{}

// LayoutFunc prepare layout data
type LayoutFunc func(H) error

// Context Context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Params   map[string]string
	Payload  H
	decoder  *form.Decoder
	validate *validator.Validate
	code     int
	index    uint8
	handlers []HandlerFunc
	render   *render.Render
}

// Bind bind json request
func (p *Context) Bind(f interface{}) error {
	if p.Request.Body == nil {
		return errors.New("nil request body")
	}
	if e := json.NewDecoder(p.Request.Body).Decode(f); e != nil {
		return e
	}
	if e := p.decoder.Decode(f, p.Request.Form); e != nil {
		return e
	}
	return p.validate.Struct(f)
}

// Form buind form request
func (p *Context) Form(f interface{}) error {
	if e := p.Request.ParseForm(); e != nil {
		return e
	}
	if e := p.decoder.Decode(f, p.Request.Form); e != nil {
		return e
	}
	return p.validate.Struct(f)
}

// ClientIP http client ip
func (p *Context) ClientIP() string {
	ip := p.Request.Header.Get("X-Forwarded-For")
	if idx := strings.IndexByte(ip, ','); idx >= 0 {
		ip = ip[0:idx]
	}
	ip = strings.TrimSpace(ip)
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(p.Request.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip = p.Request.Header.Get("X-Appengine-Remote-Addr"); ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(p.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// Next Next should be used only inside middleware. It executes the pending handlers in the chain inside the calling handler.
func (p *Context) Next() {
	if p.index >= uint8(len(p.handlers)) {
		return
	}
	h := p.handlers[p.index]
	p.index++
	log.Debugf("call %s", FuncName(h))
	h(p)
}

// Header set write header
func (p *Context) Header(k, v string) {
	p.Writer.Header().Set(k, v)
}

// Abort Abort prevents pending handlers from being called. Note that this will not stop the current handler.
func (p *Context) Abort(code int, err error) {
	http.Error(p.Writer, err.Error(), code)
	p.index = math.MaxUint8
	p.code = code
	log.Error(err)
}

// JSON render json
func (p *Context) JSON(c int, v interface{}) {
	p.render.JSON(p.Writer, c, v)
}

// XML render xml
func (p *Context) XML(c int, v interface{}) {
	p.render.XML(p.Writer, c, v)
}

// HTML render html
func (p *Context) HTML(c int, l, n string, v H) {
	p.render.HTML(p.Writer, c, n, v, render.HTMLOptions{Layout: l})
}

// TEXT render text
func (p *Context) TEXT(c int, n string, v interface{}) {
	tpl, err := template.ParseFiles(path.Join("templates", n))
	var buf bytes.Buffer
	if err == nil {
		err = tpl.Execute(&buf, v)
	}
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
		return
	}

	p.render.Text(p.Writer, c, buf.String())
}
