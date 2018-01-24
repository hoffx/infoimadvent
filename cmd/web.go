package cmd

import (
	"html/template"
	"log"

	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/parser"
	"github.com/hoffx/infoimadvent/routes"
	"github.com/hoffx/infoimadvent/services"
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
var rStorer storage.RelationStorer

func runWeb(ctx *cli.Context) {
	// load config
	config.Load(ctx.GlobalString("config"))

	// init storers
	if !qStorer.Active || !uStorer.Active || !rStorer.Active {
		var err error
		uStorer, qStorer, rStorer, err = storage.InitStorers()
		if err != nil {
			log.Fatal(err)
		}
	}

	// write admin to db
	user, err := uStorer.Get(map[string]interface{}{"email": config.Config.Auth.AdminMail})
	if err != nil {
		log.Fatal(err)
	} else if user.Email == "" {
		err = uStorer.Create(storage.User{config.Config.Auth.AdminMail, config.Config.Auth.AdminHash, config.Config.Grades.Max, true, true, "", true, make([]int, 24), 0, true})
		if err != nil {
			log.Fatal(err)
		}
	}

	// set up score-calculation service

	s := services.NewDBStorage(&uStorer, &qStorer, &rStorer)
	s.SetupRoutines()

	// set up web service

	m := initMacaron()

	m.Run(config.Config.Server.Ip, config.Config.Server.Port)
}

func initMacaron() *macaron.Macaron {
	m := macaron.New()

	mp := make(map[interface{}]interface{})
	mp["user"] = storage.User{}
	session.EncodeGob(mp)

	if config.Config.Server.DevMode == true {
		macaron.Env = macaron.DEV
	} else {
		macaron.Env = macaron.PROD
	}

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
		Funcs: []template.FuncMap{map[string]interface{}{
			"add": parser.Add,
		}},
	}))
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: config.Config.Sessioner.StoragePath,
	}))

	m.Map(&qStorer)
	m.Map(&uStorer)
	m.Map(&rStorer)

	m.Get("/", routes.Home)
	m.Route("/login", "GET,POST", routes.Login)
	m.Get("/about", routes.About)
	m.Route("/register", "GET,POST", routes.Register)
	m.Get("/confirm", routes.Confirm)
	m.Post("/restore", routes.Restore)

	m.Group("", func() {
		m.Route("/upload", "GET,POST", routes.Upload)
		m.Get("/overview", routes.Overview)
	}, routes.RequireAdmin)

	m.Group("", func() {
		m.Get("/logout", routes.Logout)
		m.Route("/account", "GET,POST", routes.Account)
	}, routes.Protect)

	m.Group("", func() {
		m.Get("/calendar", routes.Calendar)
		m.Group("/day", func() {
			m.Get("/", routes.Current)
			m.Get("/:day", routes.Day)
		})
	}, routes.PublicReady, routes.Protect)

	return m
}
