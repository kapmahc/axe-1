package axe

import (
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-playground/form"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// HandlerFunc http handler
type HandlerFunc func(*Context)

// H hash data
type H map[string]interface{}

// Context Context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Params   map[string]string
	Payload  H
	decoder  *form.Decoder
	validate *validator.Validate
	code     int
	index    int8
	handlers []HandlerFunc
}

// Bind bind json form
func (p *Context) Bind(f interface{}) error {
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
	p.index++
	s := int8(len(p.handlers))
	for ; p.index < s; p.index++ {
		h := p.handlers[p.index]
		log.Debugf("call %s", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
		h(p)
	}
}

// Abort Abort prevents pending handlers from being called. Note that this will not stop the current handler.
func (p *Context) Abort(code int, err error) {
	http.Error(p.Writer, err.Error(), code)
	p.index = math.MaxInt8
	p.code = code
	log.Error(err)
}
