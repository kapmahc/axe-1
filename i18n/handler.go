package i18n

import (
	"context"
	"net/http"

	"github.com/kapmahc/axe"
	"github.com/urfave/negroni"

	"golang.org/x/text/language"
)

const (
	// LOCALE locale key
	LOCALE = axe.K("locale")
)

// ServeHTTP middleware
func (p *I18n) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

}

// Middleware detect language from http request
func (p *I18n) Middleware() (negroni.HandlerFunc, error) {
	langs, err := p.Store.Languages()
	if err != nil {
		return nil, err
	}
	var tags []language.Tag
	for _, l := range langs {
		tags = append(tags, language.Make(l))
	}
	matcher := language.NewMatcher(tags)

	return func(wrt http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		tag, _, _ := matcher.Match(language.Make(p.detect(req)))
		ctx := context.WithValue(req.Context(), LOCALE, tag.String())
		next(wrt, req.WithContext(ctx))
	}, nil
}

func (p *I18n) detect(r *http.Request) string {
	const key = "locale"
	// 1. Check URL arguments.
	if lang := r.URL.Query().Get(key); lang != "" {
		return lang
	}

	// 2. Get language information from cookies.
	if ck, er := r.Cookie(key); er == nil {
		return ck.Value
	}

	// 3. Get language information from 'Accept-Language'.
	if al := r.Header.Get("Accept-Language"); len(al) > 4 {
		return al[:5] // Only compare first 5 letters.
	}

	return ""
}
