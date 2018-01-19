package cmd

import (
	"log"
	"os"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
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
		},
	},
}

func reset(ctx *cli.Context) {
	if config.Config.DB.Name == "" {
		config.Load(ctx.GlobalString("config"))
	}
	if !uStorer.Active || !qStorer.Active {
		initStorer()
	}

	if ctx.Bool("quests") {
		err := resetQuests()
		if err != nil {
			log.Fatal(err)
		}
	}
	if ctx.Bool("users") {
		err := resetUsers()
		if err != nil {
			log.Fatal(err)
		}
	}
	if ctx.Bool("web") {
		runWeb(ctx)
	}
}

func resetUsers() (err error) {
	err = os.RemoveAll(config.Config.Sessioner.StoragePath)
	if err != nil {
		return
	}
	err = os.Mkdir(config.Config.Sessioner.StoragePath, os.ModePerm)
	if err != nil {
		return
	}
	_, err = os.Create(config.Config.Sessioner.StoragePath + "/keep.me")
	if err != nil {
		return
	}
	err = uStorer.ResetDB()
	return
}

func resetQuests() (err error) {
	err = os.RemoveAll(config.Config.FileSystem.MDStoragePath)
	if err != nil {
		return
	}
	err = os.Mkdir(config.Config.FileSystem.MDStoragePath, os.ModePerm)
	if err != nil {
		return
	}
	_, err = os.Create(config.Config.FileSystem.MDStoragePath + "/keep.me")
	if err != nil {
		return
	}
	err = os.RemoveAll(config.Config.FileSystem.AssetsStoragePath)
	if err != nil {
		return
	}
	err = os.Mkdir(config.Config.FileSystem.AssetsStoragePath, os.ModePerm)
	if err != nil {
		return
	}
	_, err = os.Create(config.Config.FileSystem.AssetsStoragePath + "/keep.me")
	if err != nil {
		return
	}
	err = qStorer.ResetDB()
	return
}

func initStorer() {
	var err error
	uStorer, err = storage.NewUserStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		log.Fatal(err)
	}
	qStorer, err = storage.NewQuestStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		log.Fatal(err)
	}
}
