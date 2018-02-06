package cmd

import (
	"log"
	"time"

	"github.com/hoffx/infoimadvent/storage"
	"github.com/urfave/cli"
)

var Calc = cli.Command{
	Name:   "calc",
	Usage:  "calculates the user's scores for the day before execution time",
	Action: calc,
}

func calc(ctx *cli.Context) {
	setupSystem(ctx.GlobalString("config"))

	calcOperation()
}

func calcOperation() {
	users, err := uStorer.GetAll(map[string]interface{}{})
	if err != nil {
		log.Println(err)
		return
	}

	_, m, d := time.Now().Date()

	// TODO: change back to december after testing
	if m != time.February {
		return
	}

	for _, u := range users {
		// should be executed soon after the calculated day has ended
		// (e.g. next day at 1am) => -1 for last day + -1 for slice index-shift
		quest, err := dStorer.Get(map[string]interface{}{"day": d, "grade": u.Grade, "type": storage.Quest})
		if err != nil {
			log.Println(err)
			continue
		}

		if u.Days[d-2] == quest.Solution {
			u.Score += storage.Right
		} else if u.Days[d-2] == storage.None {
			u.Score += storage.Missing
		} else {
			u.Score += storage.Wrong
		}
	}
}
