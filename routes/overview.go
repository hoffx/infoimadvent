package routes

import (
	"log"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Overview(ctx *macaron.Context, log *log.Logger, qStorer *storage.QuestStorer) {
	var days [24]map[int]string

	for day := 1; day <= 24; day++ {
		days[day-1] = make(map[int]string, 0)
		for grade := config.Config.Grades.Min; grade <= config.Config.Grades.Max; grade++ {
			q, err := qStorer.Get(map[string]interface{}{"grade": grade, "day": day})
			if err != nil {
				days[day-1][int(grade)] = "!/-"
			} else if q.Path == "" {
				days[day-1][int(grade)] = "---"
			} else {
				days[day-1][int(grade)] = "+++"
			}
		}
	}
	ctx.Data["Days"] = days
	ctx.HTML(200, "overview")
}
