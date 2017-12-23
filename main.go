package main

import (
	"os"

	"github.com/hoffx/infoimadvent/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "iia"
	app.Usage = "A dynamic website for students of all ages to have fun before Christmas called IT during Advent"
	app.Author = "Hoff Industires"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.ini",
			Usage: "Configuration file path",
		},
	}
	app.Commands = []cli.Command{cmd.Web}
	app.Run(os.Args)
}
