package routes

import (
	"html/template"
	"io/ioutil"
	"log"
	"path"

	"github.com/hoffx/infoimadvent/parser"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func About(ctx *macaron.Context, qStorer *storage.QuestStorer, log *log.Logger) {

	quest, err := qStorer.Get(map[string]interface{}{"is_about": true})
	if err != nil {
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(quest.Path)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrFS))
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}

	name := Name(path.Base(quest.Path))
	html, err := parser.ParseAndProcess(data, []func(*string) error{name.parseUrls})
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}
	ctx.Data["Text"] = template.HTML(html)

	ctx.HTML(200, "about")
}
