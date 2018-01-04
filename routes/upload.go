package routes

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Upload(ctx *macaron.Context, log *log.Logger, qStorer *storage.QuestStorer) {
	defer ctx.HTML(200, "upload")

	ctx.Data["MinYL"] = config.Config.Grades.Min
	ctx.Data["MaxYL"] = config.Config.Grades.Max

	if ctx.Req.Method == "GET" {
		return
	} else {
		fPw := ctx.Req.FormValue("pw")
		fMinGrade, err := strconv.Atoi(ctx.Req.FormValue("mingrade"))
		if err != nil {
			ctx.Data["Error"] = ErrIllegalInput
			return
		}
		fMaxGrade, err := strconv.Atoi(ctx.Req.FormValue("maxgrade"))
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
		ctx.Data["MinGrade"] = fMinGrade
		ctx.Data["MaxGrade"] = fMaxGrade
		ctx.Data[fSolution] = true
		ctx.Data["Markdown"] = fMd

		solution, err := solutionToInt(fSolution)
		if err != nil {
			ctx.Data["Error"] = ErrIllegalInput
			return
		}

		f, err := ioutil.TempFile(config.Config.FileSystem.StoragePath, "quest")
		if err != nil {
			ctx.Data["Error"] = ErrFS
			log.Println(err)
			return
		}
		defer f.Close()

		_, err = f.WriteString(fMd)
		if err != nil {
			ctx.Data["Error"] = ErrFS
			log.Println(err)
			return
		}

		for i := fMinGrade; i <= fMaxGrade; i++ {
			quest := storage.Quest{f.Name(), i, fDay, solution}
			oldQ, err := qStorer.Get(map[string]interface{}{"grade": i, "day": fDay})
			if err != nil {
				ctx.Data["Error"] = ErrDB
				log.Println(err)
				return
			}
			if oldQ.Path == "" {
				err = qStorer.Create(quest)
				if err != nil {
					ctx.Data["Error"] = ErrDB
					log.Println(err)
					return
				}
			} else {
				err = qStorer.Put(quest)
				if err != nil {
					ctx.Data["Error"] = ErrDB
					log.Println(err)
					return
				}
			}
		}

	}
}

func solutionToInt(sol string) (solution int, err error) {
	switch sol {
	case "":
		solution = storage.None
	case "A":
		solution = storage.A
	case "B":
		solution = storage.B
	case "C":
		solution = storage.C
	case "D":
		solution = storage.D
	default:
		err = errors.New(ErrIllegalInput)
	}
	return
}
