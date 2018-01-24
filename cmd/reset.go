package cmd

import (
	"log"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
)

var Reset = cli.Command{
	Name:   "reset",
	Usage:  "resets database and session-storage",
	Action: reset,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "users, u",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "quests, q",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "web, w",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "all, a",
			Hidden: false,
		},
	},
}

func reset(ctx *cli.Context) {
	if config.Config.DB.Name == "" {
		config.Load(ctx.GlobalString("config"))
	}
	if !uStorer.Active || !qStorer.Active || !rStorer.Active {
		var err error
		uStorer, qStorer, rStorer, err = storage.InitStorers()
		if err != nil {
			log.Fatal(err)
		}
	}

	if ctx.Bool("quests") || ctx.Bool("all") {
		err := storage.ResetQuests(&qStorer)
		if err != nil {
			log.Fatal(err)
		}
	}
	if ctx.Bool("users") || ctx.Bool("all") {
		err := storage.ResetUsers(&uStorer, &rStorer)
		if err != nil {
			log.Fatal(err)
		}
	}
	if ctx.Bool("web") {
		runWeb(ctx)
	}
}
