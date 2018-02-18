package main

import (
	"os"

	"github.com/hoffx/infoimadvent/cmd"
	"github.com/urfave/cli"
)

// Version holds the version-string. It is set during the
// build process via the "make release" command.
var Version string

// GitCommit holds the comment-id of the latest commit at
// build time. It is set during the build process via
// the "make release" command.
var GitCommit string

// DefaultConfigPath holds the default path to the configuration
// file. It is set during the build process via the "make
// release" command.
var DefaultConfigPath string

func main() {
	if DefaultConfigPath == "" {
		DefaultConfigPath = "config.ini"
	}
	app := cli.NewApp()
	app.Name = "iia"
	app.Version = Version + " (" + GitCommit + ")"
	app.Usage = "A dynamic website for students of all ages to have fun before Christmas called IT during Advent"
	app.Author = "Hoff Industires"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: DefaultConfigPath,
			Usage: "Configuration file path",
		},
	}
	app.Commands = []cli.Command{cmd.Web, cmd.Reset, cmd.Calc}
	app.Run(os.Args)
}
