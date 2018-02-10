package routes

import (
	"log"
	"strings"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Certificate(ctx *macaron.Context, sess session.Store) {

	if !certificateReady() {
		ctx.Error(503, ctx.Tr(ErrAdventNotOver))
		return
	}

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

	if err = pdfg.Create(); err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	if err = pdfg.WriteFile("pdf/certificate.pdf"); err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/account", 500)
		return
	}

	ctx.ServeFile("pdf/certificate.pdf")
}

func certificateReady() bool {
	_, month, day := time.Now().Date()

	return month != time.February || day > 24
}
