package cmd

import (
	"log"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
)

var Reset = cli.Command{
	Name:   "reset",
	Usage:  "resets database",
	Action: resetDB,
}

func resetDB(ctx *cli.Context) {
	if config.Config.DB.Name == "" {
		config.Load(ctx.GlobalString("config"))
	}
	if !storer.Active {
		initStorer()
	}
	err := storer.ResetDB()
	if err != nil {
		log.Println(err)
	}

	if ctx.Args().First() == "web" {
		runWeb(ctx)
	}
}

func initStorer() {
	var err error
	storer, err = storage.NewStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		log.Fatal(err)
	}
}
