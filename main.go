package main

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kapmahc/axe/app"
	_ "github.com/kapmahc/axe/plugins/erp"
	_ "github.com/kapmahc/axe/plugins/forum"
	_ "github.com/kapmahc/axe/plugins/mall"
	_ "github.com/kapmahc/axe/plugins/nut"
	_ "github.com/kapmahc/axe/plugins/ops/mail"
	_ "github.com/kapmahc/axe/plugins/ops/vpn"
	_ "github.com/kapmahc/axe/plugins/pos"
	_ "github.com/kapmahc/axe/plugins/reading"
	_ "github.com/kapmahc/axe/plugins/survey"
)

func main() {
	if err := app.Main(); err != nil {
		log.Fatal(err)
	}
}