package cmd

import (
	"log"

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
			Name:   "docs, d",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "web, w",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "all, a",
			Hidden: false,
		}, cli.BoolFlag{
			Name:   "standard, s",
			Hidden: false,
		},
	},
}

func reset(ctx *cli.Context) {
	setupSystem(ctx.GlobalString("config"))

	if ctx.Bool("standard") || (!ctx.Bool("docs") && !ctx.Bool("all") && !ctx.Bool("users")) {
		standardReset()
	} else if ctx.Bool("docs") || ctx.Bool("all") {
		err := storage.ResetDocuments(&dStorer, false)
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

func standardReset() {
	err := storage.ResetDocuments(&dStorer, true)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.ResetUsers(&uStorer, &rStorer)
	if err != nil {
		log.Fatal(err)
	}
}
