package base

import (
	"fmt"
	"log/syslog"
	"reflect"

	"github.com/facebookgo/inject"
	"github.com/kapmahc/axe"
	"github.com/kapmahc/axe/job"
	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// IsProduction production mode ?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}

// Name server name
func Name() string {
	return viper.GetString("server.name")
}

// Home home url
func Home() string {
	host := Name()
	if IsProduction() {
		if viper.GetBool("server.ssl") {
			return "https://" + host
		}
		return "http://" + host
	}
	return fmt.Sprintf("http://localhost:%d", viper.GetInt("server.port"))
}

type injectLogger struct {
}

func (p *injectLogger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Config load config first
func Config(f cli.ActionFunc) cli.ActionFunc {

	viper.SetEnvPrefix(reflect.TypeOf(axe.Main).String())
	viper.BindEnv("env")
	viper.SetDefault("env", "development")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	return func(c *cli.Context) error {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		if IsProduction() {
			// ----------
			log.SetLevel(log.InfoLevel)
			if wrt, err := syslog.New(syslog.LOG_INFO, Name()); err == nil {
				log.AddHook(&logrus_syslog.SyslogHook{Writer: wrt})
			} else {
				log.Error(err)
			}
		} else {
			log.SetLevel(log.DebugLevel)
		}

		log.Infof("read config from %s", viper.ConfigFileUsed())

		return f(c)
	}
}

// Inject inject objects first
func Inject(f func(*cli.Context, *inject.Graph) error) cli.ActionFunc {
	return Config(func(c *cli.Context) error {
		gh := inject.Graph{Logger: &injectLogger{}}

		srv := job.New()
		if err := gh.Provide(
			&inject.Object{Value: srv},
		); err != nil {
			return err
		}

		if err := axe.Walk(func(p axe.Plugin) error {
			if err := p.Open(&gh); err != nil {
				return err
			}
			return gh.Provide(&inject.Object{Value: p})
		}); err != nil {
			return err
		}

		if err := gh.Populate(); err != nil {
			return err
		}
		// ------------

		axe.Walk(func(p axe.Plugin) error {
			for k, v := range p.Workers() {
				srv.Register(k, v)
			}
			return nil
		})
		// ------------
		return f(c, &gh)
	})
}
