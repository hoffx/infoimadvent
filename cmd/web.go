package cmd

import (
	"html/template"
	"image/color"
	"log"
	"strconv"

	"github.com/go-macaron/cache"
	"github.com/robfig/cron"
	"github.com/theMomax/captcha"

	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/parser"
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
var dStorer storage.DocumentStorer
var rStorer storage.RelationStorer

func runWeb(ctx *cli.Context) {
	setupSystem(ctx.GlobalString("config"))

	startCronJobs()

	// write admin to db
	user, err := uStorer.Get(map[string]interface{}{"email": config.Config.Auth.AdminMail})
	if err != nil {
		log.Fatal(err)
	} else if user.Email == "" {
		err = uStorer.Create(storage.User{config.Config.Auth.AdminMail, config.Config.Auth.AdminHash, config.Config.Grades.Max, true, true, "", true, make([]int, 24), 0, true, "en-US"})
		if err != nil {
			log.Fatal(err)
		}
	}

	// set up web service

	m := initMacaron()

	m.Run(config.Config.Server.Ip, config.Config.Server.Port)
}

func startCronJobs() {
	c := cron.New()
	c.AddFunc("0 0 1 "+strconv.Itoa(config.Config.Server.ResetMonth)+" *", standardReset)
	// TODO: change back to december after testing
	c.AddFunc("0 2 2-25 2 *", calcOperation)
	c.Start()
}

func setupSystem(configpath string) {
	// check if config is already loaded
	if config.Config.DB.Name == "" {
		config.Load(configpath)
	}
	// check if storage has been activated
	if !uStorer.Active || !dStorer.Active || !rStorer.Active {
		var err error
		uStorer, dStorer, rStorer, err = storage.InitStorers()
		if err != nil {
			log.Fatal(err)
		}
	}
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
	m.Use(cache.Cacher())
	m.Use(captcha.Captchaer(
		captcha.Options{
			ColorPalette: buildColorPalette(),
		},
	))

	m.Map(&dStorer)
	m.Map(&uStorer)
	m.Map(&rStorer)

	m.Get("/", routes.Home)
	m.Route("/login", "GET,POST", routes.Login)
	m.Get("/about", routes.About)
	m.Get("/tos", routes.ToS)
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

func buildColorPalette() (cp color.Palette) {
	return color.Palette{color.RGBA{196, 187, 69, 255}, color.RGBA{65, 144, 42, 255}, color.RGBA{210, 78, 76, 255}, color.RGBA{210, 210, 210, 255}}
}
