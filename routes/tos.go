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

func ToS(ctx *macaron.Context, dStorer *storage.DocumentStorer, log *log.Logger) {

	doc, err := dStorer.Get(map[string]interface{}{"type": storage.ToS})
	if err != nil {
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(doc.Path)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrFS))
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}

	name := Name(path.Base(doc.Path))
	html, err := parser.ParseAndProcess(data, []func(*string) error{name.parseUrls})
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}
	ctx.Data["Text"] = template.HTML(html)

	ctx.HTML(200, "tos")
}
