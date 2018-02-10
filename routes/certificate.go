package routes

import (
	"log"
	"strings"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Certificate(ctx *macaron.Context, sess session.Store) {
	user := sess.Get("user").(storage.User) // protected therefore user must exist

	ctx.Data["Email"] = user.Email
	ctx.Data["Score"] = user.Score

	htmlBody, err := ctx.HTMLString("certificate", ctx.Data)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	page := wkhtmltopdf.NewPageReader(strings.NewReader(htmlBody))
	page.PrintMediaType.Set(true)
	page.DisableSmartShrinking.Set(true)

	pdfg.AddPage(page)

	err = pdfg.Create()
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	err = pdfg.WriteFile("pdf/certificate.pdf")
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	ctx.Redirect("/account?mode=certificate", 302)
}
