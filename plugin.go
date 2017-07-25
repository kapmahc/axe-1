package axe

import (
	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/kapmahc/axe/job"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
	"golang.org/x/tools/blog/atom"
)

// Plugin plugin
type Plugin interface {
	Open(*inject.Graph) error
	Mount(*mux.Router, *negroni.Negroni)
	Console() []cli.Command
	Workers() map[string]job.Handler
	Rss(lang string) ([]*atom.Entry, error)
	Sitemap() ([]stm.URL, error)
}

var plugins []Plugin

// Register register plugin
func Register(args ...Plugin) {
	plugins = append(plugins, args...)
}

// Walk walk plugins
func Walk(fn func(Plugin) error) error {
	for _, p := range plugins {
		if e := fn(p); e != nil {
			return e
		}
	}
	return nil
}
