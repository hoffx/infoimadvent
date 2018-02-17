package cmd

import (
	"log"

	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
)

// Reset holds the cli command for all available reset methods
var Reset = cli.Command{
	Name:   "reset",
	Usage:  "resets database and session-storage",
	Action: reset,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "users, u",
			Usage: "deletes all users except admin",
		}, cli.BoolFlag{
			Name:  "docs, d",
			Usage: "deletes all documents of all types",
		}, cli.BoolFlag{
			Name:  "all, a",
			Usage: "deletes all documents of all types and all users except admin",
		}, cli.BoolFlag{
			Name:  "standard, s",
			Usage: "deletes all documents of type quest and all users except admin",
		}, cli.BoolFlag{
			Name:  "web, w",
			Usage: "starts the webserver after executing the reset",
		},
	},
}

// reset interprets the cli command "reset"
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

// standardReset deletes all documents of type quest
// and all users (except admin)
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
