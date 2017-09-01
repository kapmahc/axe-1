package vpn

import (
	"github.com/facebookgo/inject"
	"github.com/graphql-go/graphql"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/urfave/cli"
	"golang.org/x/tools/blog/atom"
)

// Plugin plugin
type Plugin struct {
}

// Graphql register graphql fields
func (p *Plugin) Graphql(q graphql.Fields, m graphql.Fields) {}

// Open init beans
func (p *Plugin) Open(*inject.Graph) error {
	return nil
}

// Console shell commands
func (p *Plugin) Console() []cli.Command {
	return []cli.Command{}
}

// Atom rss.atom
func (p *Plugin) Atom(lang string) ([]*atom.Entry, error) {
	return []*atom.Entry{}, nil
}

// Sitemap sitemap.xml.gz
func (p *Plugin) Sitemap() ([]stm.URL, error) {
	return []stm.URL{}, nil
}

// Workers background workers
func (p *Plugin) Workers() map[string]nut.JobHandler {
	return make(map[string]nut.JobHandler)
}

func init() {
	nut.Register(&Plugin{})
}
