package cmd

import (
	"github.com/go-macaron/i18n"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/routes"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
)

var Web = cli.Command{
	Name:   "web",
	Usage:  "Start webserver",
	Action: runWeb,
}

func runWeb(ctx *cli.Context) {
	config.Load(ctx.GlobalString("config"))

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static", macaron.StaticOptions{
		SkipLogging: true,
	}))
	m.Use(i18n.I18n(i18n.Options{
		Directory: "locales",
		Langs:     []string{"de-DE", "en-US"},
		Names:     []string{"Deutsch", "Englisch"},
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates",
	}))
	m.Get("/", routes.Home)
	m.Get("/calendar", routes.Calendar)
	m.Get("/register", routes.Register)
	m.Get("/login", routes.Login)
	m.Get("/current", routes.Current)
	m.Get("/about", routes.About)
	m.Run()
}
