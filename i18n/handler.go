package i18n

import (
	"math"
	"net/http"

	"github.com/kapmahc/axe"
	log "github.com/sirupsen/logrus"
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
	log.Info("allow languages ", langs)
	var tags []language.Tag
	for _, l := range langs {
		tags = append(tags, language.Make(l))
	}
	matcher := language.NewMatcher(tags)

	return func(c *axe.Context) {
		lang, write := p.detect(c.Request)
		tag, _, _ := matcher.Match(language.Make(lang))
		lang = tag.String()
		if write {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     LOCALE,
				Value:    lang,
				Secure:   false,
				HttpOnly: false,
				MaxAge:   math.MaxUint32,
			})
		}

		c.Payload[LOCALE] = lang
		c.Payload["languages"] = langs
		c.Next()
	}, nil
}

func (p *I18n) detect(r *http.Request) (string, bool) {
	// 1. Check URL arguments.
	if lang := r.URL.Query().Get(LOCALE); lang != "" {
		return lang, true
	}

	// 2. Get language information from cookies.
	if ck, er := r.Cookie(LOCALE); er == nil {
		return ck.Value, false
	}

	// 3. Get language information from 'Accept-Language'.
	if al := r.Header.Get("Accept-Language"); len(al) > 4 {
		return al[:5], true // Only compare first 5 letters.
	}

	return "", true
}
