package cmd

import (
	"log"
	"time"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
)

// Calc holds the cli command for calculating
var Calc = cli.Command{
	Name:   "calc",
	Usage:  "calculates the user's scores for the day before execution time",
	Action: calc,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "day, d",
			Usage: "calculates the scores for the given day (to compensate server-down-time)",
		},
	},
}

// calc interprets the cli command "calc"
func calc(ctx *cli.Context) {
	setupSystem(ctx.GlobalString("config"))
	d := ctx.Int("day")
	if d <= 24 && d >= 1 && d < time.Now().Day() {
		calcOperation(d)
	} else if d == 0 {
		calcLatest()
	} else {
		log.Fatal("Day makes no sense !")
	}
}

// calcLatest is a shortcut for calculation the latest finished day.
func calcLatest() {
	calcOperation(time.Now().Day() - 1)
}

// calcOperation calculates the user's scores for the given day.
func calcOperation(d int) {

	m := time.Now().Month()

	if m != config.Config.Server.Advent {
		return
	}

	users, err := uStorer.GetAll(map[string]interface{}{})
	if err != nil {
		log.Println(err)
		return
	}

	for _, u := range users {
		quest, err := dStorer.Get(map[string]interface{}{"day": d, "grade": u.Grade, "type": storage.Quest})
		if err != nil {
			log.Println(err)
			continue
		}

		if u.Days[d-1] == quest.Solution {
			u.Score += storage.Right
		} else if u.Days[d-1] == storage.None {
			u.Score += storage.Missing
		} else {
			u.Score += storage.Wrong
		}

		err = uStorer.Put(u)
		if err != nil {
			log.Println(err)
		}
	}
}
