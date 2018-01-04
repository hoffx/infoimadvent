package routes

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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
		// parse trivial form values
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

		// save trivial form values
		ctx.Data["Pw"] = fPw
		ctx.Data["Day"] = fDay
		ctx.Data["MinGrade"] = fMinGrade
		ctx.Data["MaxGrade"] = fMaxGrade
		ctx.Data[fSolution] = true

		// parse non-trivial form values

		if fPw != config.Config.Auth.AdminPassword {
			ctx.Data["Error"] = ErrWrongCredentials
			return
		}

		solution, err := solutionToInt(fSolution)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			return
		}

		// load and save form files
		fMd, _, err := ctx.Req.FormFile("md")
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			log.Println(err)
			return
		}
		defer fMd.Close()

		f, err := ioutil.TempFile(config.Config.FileSystem.MDStoragePath, "quest")
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrFS)
			log.Println(err)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, fMd)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrFS)
			log.Println(err)
			return
		}

		fAssets, _, err := ctx.Req.FormFile("assets")
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			log.Println(err)
			return
		}
		defer fAssets.Close()
		buf := new(bytes.Buffer)
		length, err := buf.ReadFrom(fAssets)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			log.Println(err)
		}

		reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), length)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			log.Println(err)
			return
		}
		err = unzipAndSave(*reader, config.Config.FileSystem.AssetsStoragePath+"/"+path.Base(f.Name()))
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrFS)
			log.Println(err)
			return
		}

		// create db entries

		for i := fMinGrade; i <= fMaxGrade; i++ {
			quest := storage.Quest{f.Name(), i, fDay, solution}
			oldQ, err := qStorer.Get(map[string]interface{}{"grade": i, "day": fDay})
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrDB)
				log.Println(err)
				return
			}
			if oldQ.Path == "" {
				err = qStorer.Create(quest)
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				}
			} else {
				err = qStorer.Put(quest)
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
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

func unzipAndSave(reader zip.Reader, target string) error {

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}