package i18n

import (
	"net/http"

	"github.com/kapmahc/axe"

	"golang.org/x/text/language"
)

const (
	// LOCALE locale key
	LOCALE = "locale"
)

// Middleware detect language from http request
func (p *I18n) Middleware() (axe.HandlerFunc, error) {
	langs, err := p.Store.Languages()
	if err != nil {
		return nil, err
	}
	var tags []language.Tag
	for _, l := range langs {
		tags = append(tags, language.Make(l))
	}
	matcher := language.NewMatcher(tags)

	return func(c *axe.Context) {
		tag, _, _ := matcher.Match(language.Make(p.detect(c.Request)))
		c.Payload[LOCALE] = tag.String()
		c.Next()
	}, nil
}

func (p *I18n) detect(r *http.Request) string {
	// 1. Check URL arguments.
	if lang := r.URL.Query().Get(LOCALE); lang != "" {
		return lang
	}

	// 2. Get language information from cookies.
	if ck, er := r.Cookie(LOCALE); er == nil {
		return ck.Value
	}

	// 3. Get language information from 'Accept-Language'.
	if al := r.Header.Get("Accept-Language"); len(al) > 4 {
		return al[:5] // Only compare first 5 letters.
	}

	return ""
}
