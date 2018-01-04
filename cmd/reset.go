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
}

func reset(ctx *cli.Context) {
	if config.Config.DB.Name == "" {
		config.Load(ctx.GlobalString("config"))
	}
	if !uStorer.Active || !qStorer.Active {
		initStorer()
	}

	var web bool
	args := append(ctx.Args().Tail(), ctx.Args().First())
	for _, a := range args {
		switch a {
		case "filesystem":
			err := os.RemoveAll(config.Config.FileSystem.StoragePath)
			if err != nil {
				log.Println(err)
			}
			err = os.Mkdir(config.Config.FileSystem.StoragePath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
			_, err = os.Create(config.Config.FileSystem.StoragePath + "/keep.me")
			if err != nil {
				log.Println(err)
			}
		case "sessions":
			err := os.RemoveAll(config.Config.Sessioner.StoragePath)
			if err != nil {
				log.Println(err)
			}
			err = os.Mkdir(config.Config.Sessioner.StoragePath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
			_, err = os.Create(config.Config.Sessioner.StoragePath + "/keep.me")
			if err != nil {
				log.Println(err)
			}
		case "db":
			err := uStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
			err = qStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
		case "userDB":
			err := uStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
		case "questDB":
			err := qStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
		case "all":
			err := os.RemoveAll(config.Config.FileSystem.StoragePath)
			if err != nil {
				log.Println(err)
			}
			err = os.Mkdir(config.Config.FileSystem.StoragePath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
			_, err = os.Create(config.Config.FileSystem.StoragePath + "/keep.me")
			if err != nil {
				log.Println(err)
			}
			err = os.RemoveAll(config.Config.Sessioner.StoragePath)
			if err != nil {
				log.Println(err)
			}
			err = os.Mkdir(config.Config.Sessioner.StoragePath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
			_, err = os.Create(config.Config.Sessioner.StoragePath + "/keep.me")
			if err != nil {
				log.Println(err)
			}
			err = uStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
			err = qStorer.ResetDB()
			if err != nil {
				log.Println(err)
			}
		case "web":
			web = true
		}
	}

	if web {
		runWeb(ctx)
	}
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
