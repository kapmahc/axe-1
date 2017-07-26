package axe

import "net/http"

// Context http context
type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request: r,
		Writer:  w,
	}
}
