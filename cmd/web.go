package cmd

import (
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/routes"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
)

var Web = cli.Command{
	Name:   "web",
	Usage:  "Start webserver",
	Action: runWeb,
}

var uStorer storage.UserStorer
var qStorer storage.QuestStorer

func runWeb(ctx *cli.Context) {
	config.Load(ctx.GlobalString("config"))

	if config.Config.Server.DevMode == true {
		macaron.Env = macaron.DEV
	} else {
		macaron.Env = macaron.PROD
	}

	m := macaron.New()

	mp := make(map[interface{}]interface{})
	mp["user"] = storage.User{}
	session.EncodeGob(mp)

	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static", macaron.StaticOptions{
		SkipLogging: true,
	}))
	m.Use(macaron.Static(config.Config.FileSystem.AssetsStoragePath, macaron.StaticOptions{
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
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: config.Config.Sessioner.StoragePath,
	}))

	if !qStorer.Active || !uStorer.Active {
		initStorer()
	}

	m.Map(&qStorer)
	m.Map(&uStorer)

	m.Get("/", routes.Home)

	m.Route("/register", "GET,POST", routes.Register)
	m.Route("/login", "GET,POST", routes.Login)
	m.Get("/about", routes.About)
	m.Get("/confirm", routes.Confirm)
	m.Post("/restore", routes.Restore)
	m.Route("/upload", "GET,POST", routes.Upload)

	m.Group("", func() {
		m.Route("/account", "GET,POST", routes.Account)
		m.Get("/logout", routes.Logout)
		m.Get("/calendar", routes.Calendar)
		m.Group("/day", func() {
			m.Get("/", routes.Current)
			m.Get("/:day", routes.Day)
		})
	}, routes.Protect)

	m.Run(config.Config.Server.Ip, config.Config.Server.Port)
}
