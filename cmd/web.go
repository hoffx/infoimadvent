package cmd

import (
	"html/template"
	"image/color"
	"log"
	"strconv"
	"time"

	"github.com/go-macaron/cache"
	"github.com/jung-kurt/gofpdf"
	"github.com/robfig/cron"
	"github.com/theMomax/captcha"

	"github.com/go-macaron/gzip"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/parser"
	"github.com/hoffx/infoimadvent/routes"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
)

// Web holds the cli command for starting the webserver
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

	if config.Config.Server.DevMode {
		config.Config.Server.Advent = time.Now().Month()
	} else {
		config.Config.Server.Advent = time.December
	}

	startCronJobs()

	// write admin to db
	user, err := uStorer.Get(map[string]interface{}{"email": config.Config.Auth.AdminMail})
	if err != nil {
		log.Fatal(err)
	} else if user.Email == "" {
		err = uStorer.Create(storage.User{
			Email:             config.Config.Auth.AdminMail,
			Hash:              config.Config.Auth.AdminHash,
			Grade:             config.Config.Grades.Max,
			Active:            true,
			Confirmed:         true,
			ConfirmationToken: "",
			Teacher:           true,
			Days:              make([]int, 24),
			Score:             0,
			IsAdmin:           true,
			Lang:              "en-US",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// generate font files for PDF generations
	err = gofpdf.MakeFont("static/fonts/zillaslab.ttf", "static/fonts/cp1252.map", "static/fonts/", nil, true)
	if err != nil {
		log.Fatal(err)
	}

	// set up web service

	m := initMacaron()

	m.Run(config.Config.Server.IP, config.Config.Server.Port)
}

func startCronJobs() {
	c := cron.New()
	c.AddFunc("0 0 1 "+strconv.Itoa(config.Config.Server.ResetMonth)+" *", standardReset)
	// TODO: change back to december after testing
	c.AddFunc("0 2 2-25 "+strconv.Itoa(int(config.Config.Server.Advent))+" *", calcLatest)
	c.Start()
}

func setupSystem(configpath string) {
	// load config
	config.Load(configpath)
	// init storage
	var err error
	uStorer, dStorer, rStorer, err = storage.InitStorers()
	if err != nil {
		log.Fatal(err)
	}

}

func initMacaron() *macaron.Macaron {
	m := macaron.New()

	mp := make(map[interface{}]interface{})
	mp["user"] = storage.User{}
	session.EncodeGob(mp)

	if config.Config.Server.DevMode {
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
	m.Use(gzip.Gziper())

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
		m.Get("/certificate", routes.Certificate)
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
