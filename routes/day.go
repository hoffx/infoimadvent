package routes

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/parser"
	"github.com/hoffx/infoimadvent/storage"

	macaron "gopkg.in/macaron.v1"
)

type Name string

func Day(ctx *macaron.Context, log *log.Logger, dStorer *storage.DocumentStorer, sess session.Store, uStorer *storage.UserStorer) {
	num := ctx.ParamsInt("day")
	if num < 1 || num > 24 {
		ctx.Error(404, ctx.Tr(ErrIllegalDate))
		ctx.Redirect("/calendar", 308)
		return
	}
	_, m, d := time.Now().Date()
	// TODO: change back to december after testing
	if m != time.February || d < num {
		ctx.Error(406, ctx.Tr(ErrIllegalDate))
		ctx.Redirect("/calendar", 406)
		return
	}

	user := sess.Get("user").(storage.User)

	if tip := ctx.Req.URL.Query()["tip"]; len(tip) > 0 && tip[0] != "" {
		intTip, err := solutionToInt(tip[0])
		if err != nil {
			ctx.Error(406, ctx.Tr(ErrIllegalInput))
			ctx.Redirect("/day/"+strconv.Itoa(num), 406)
			return
		}
		user.Days[num-1] = intTip

		err = uStorer.Put(user)
		if err != nil {
			ctx.Error(500, ctx.Tr(ErrDB))
			ctx.Redirect("/", 500)
			return
		}

		err = sess.Set("user", user)
		if err != nil {
			ctx.Error(500, ctx.Tr(ErrUnexpected))
			ctx.Redirect("/", 500)
			return
		}
	}

	doc, err := dStorer.Get(map[string]interface{}{"day": num, "grade": user.Grade, "type": storage.Quest})
	if err != nil {
		ctx.Redirect("/calendar", 500)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(doc.Path)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrFS))
		ctx.Redirect("/calendar", 500)
		log.Println(err)
		return
	}

	name := Name(path.Base(doc.Path))
	html, err := parser.ParseAndProcess(data, []func(*string) error{name.parseUrls})
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		ctx.Redirect("/calendar", 500)
		log.Println(err)
		return
	}

	ctx.Data["Text"] = template.HTML(html)
	tipString, _ := solutionToString(user.Days[num-1])
	ctx.Data["Tip"+tipString] = true
	ctx.Data["Day"] = num

	if d == num {
		ctx.Data["Current"] = true
	}

	if d > num {
		solString, _ := solutionToString(doc.Solution)
		ctx.Data["Solution"+solString] = true
	}

	ctx.HTML(200, "day")
}

// adds the files name to all assets paths
func (n Name) parseUrls(s *string) (err error) {
	*s = strings.Replace(*s, `src="http`, `\\\src="http///`, -1)
	*s = strings.Replace(*s, `src="`, `src="/`+string(n)+"/", -1)
	*s = strings.Replace(*s, `\\\src="http///`, `src="http`, -1)
	return nil
}

func solutionToString(sol int) (solution string, err error) {
	switch sol {
	case storage.None:
		solution = ""
	case storage.A:
		solution = "A"
	case storage.B:
		solution = "B"
	case storage.C:
		solution = "C"
	case storage.D:
		solution = "D"
	default:
		err = errors.New(ErrIllegalInput)
	}
	return
}
