package nut

import (
	"github.com/facebookgo/inject"
	"github.com/graphql-go/graphql"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"golang.org/x/tools/blog/atom"
)

// Plugin plugin
type Plugin interface {
	Graphql(q graphql.Fields, m graphql.Fields)
	Open(*inject.Graph) error
	Console() []cli.Command
	Atom(lang string) ([]*atom.Entry, error)
	Sitemap() ([]stm.URL, error)
	Workers() map[string]JobHandler
}

var plugins []Plugin

// Register register plugins
func Register(args ...Plugin) {
	plugins = append(plugins, args...)
}

// Walk walk plugins
func Walk(f func(Plugin) error) error {
	for _, p := range plugins {
		if err := f(p); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	viper.SetDefault("s3", map[string]interface{}{
		"endpoint":   "http://localhost:9000",
		"access-key": "guest",
		"secret-key": "change-me",
	})

	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   8,
	})

	viper.SetDefault("rabbitmq", map[string]interface{}{
		"user":     "guest",
		"password": "guest",
		"host":     "localhost",
		"port":     "5672",
		"virtual":  "axe-dev",
	})

	viper.SetDefault("database", map[string]interface{}{
		"driver": "postgres",
		"args": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"user":     "postgres",
			"password": "",
			"dbname":   "axe_dev",
			"sslmode":  "disable",
		},
		"pool": map[string]int{
			"max_open": 180,
			"max_idle": 6,
		},
	})

	viper.SetDefault("server", map[string]interface{}{
		"port":     8080,
		"name":     "change-me.com",
		"frontend": []string{"http://localhost:3000"},
		"backend":  "http://localhost:8080",
	})

	viper.SetDefault("secrets", map[string]interface{}{
		"jwt":  Random(32),
		"aes":  Random(32),
		"hmac": Random(32),
	})

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
	})

}

// NPlugin nut plugin
type NPlugin struct {
}

// Graphql register graphql fields
func (p *NPlugin) Graphql(q graphql.Fields, m graphql.Fields) {}

// Open init beans
func (p *NPlugin) Open(*inject.Graph) error {
	return nil
}

// Console shell commands
func (p *NPlugin) Console() []cli.Command {
	return []cli.Command{}
}

// Atom rss.atom
func (p *NPlugin) Atom(lang string) ([]*atom.Entry, error) {
	return []*atom.Entry{}, nil
}

// Sitemap sitemap.xml.gz
func (p *NPlugin) Sitemap() ([]stm.URL, error) {
	return []stm.URL{}, nil
}

// Workers background workers
func (p *NPlugin) Workers() map[string]JobHandler {
	return make(map[string]JobHandler)
}

func init() {
	Register(&NPlugin{})
}
