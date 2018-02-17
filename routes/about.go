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

// About handles the route "/about"
func About(ctx *macaron.Context, dStorer *storage.DocumentStorer, log *log.Logger) {
	// get about document's db entry
	doc, err := dStorer.Get(map[string]interface{}{"type": storage.About})
	if err != nil {
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}
	// get about md file
	data, err := ioutil.ReadFile(doc.Path)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrFS))
		ctx.Redirect("/", 500)
		log.Println(err)
		return
	}
	// render about md file to html
	name := Name(path.Base(doc.Path))
	// edit urls after parsing so that they fit the server's file-structure
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
