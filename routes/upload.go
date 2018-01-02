package routes

import (
	"strconv"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Upload(ctx *macaron.Context, qStorer *storage.QuestStorer) {
	defer ctx.HTML(200, "upload")

	ctx.Data["MinYL"] = config.Config.Grades.Min
	ctx.Data["MaxYL"] = config.Config.Grades.Max

	if ctx.Req.Method == "GET" {
		return
	} else {
		fPw := ctx.Req.FormValue("pw")
		fGrade, err := strconv.Atoi(ctx.Req.FormValue("grade"))
		if err != nil {
			ctx.Data["Error"] = ErrIllegalInput
			return
		}
		fDay, err := strconv.Atoi(ctx.Req.FormValue("day"))
		if err != nil {
			ctx.Data["Error"] = ErrIllegalInput
			return
		}
		fSolution := ctx.Req.FormValue("solution")
		fMd := ctx.Req.FormValue("md")

		if fPw != config.Config.Auth.AdminPassword {
			ctx.Data["Error"] = ErrWrongCredentials
			return
		}

		ctx.Data["Pw"] = fPw
		ctx.Data["Day"] = fDay
		ctx.Data["Grade"] = fGrade
		ctx.Data[fSolution] = true
		ctx.Data["Markdown"] = fMd

		var solution int
		switch fSolution {
		case "A":
			solution = storage.A
		case "B":
			solution = storage.B
		case "C":
			solution = storage.C
		case "D":
			solution = storage.D
		default:
			ctx.Data["Error"] = ErrIllegalInput
			return
		}

		// TODO: store file
		var path string = fMd

		quest := storage.Quest{path, fGrade, fDay, solution}
		oldQ, err := qStorer.Get(map[string]interface{}{"grade": fGrade, "day": fDay})
		if err != nil {
			ctx.Data["Error"] = ErrDB
			return
		}
		if oldQ.Path == "" {
			err = qStorer.Create(quest)
			if err != nil {
				ctx.Data["Error"] = ErrDB
				return
			}
		} else {
			err = qStorer.Put(quest)
			if err != nil {
				ctx.Data["Error"] = ErrDB
				return
			}
		}

	}
}
