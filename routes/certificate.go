package routes

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"github.com/jung-kurt/gofpdf"
	macaron "gopkg.in/macaron.v1"
)

func Certificate(ctx *macaron.Context, sess session.Store, log *log.Logger) {

	if !certificateReady() {
		ctx.Error(503, ctx.Tr(ErrAdventNotOver))
		return
	}

	user := sess.Get("user").(storage.User) // protected therefore user must exist

	ctx.Data["Email"] = user.Email
	ctx.Data["Score"] = user.Score

	file := new(bytes.Buffer)
	err := generateCertificate(file, user, ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		ctx.Redirect("/account", 500)
		return
	}

	ctx.Resp.Header().Add("content-disposition", "attachment; filename="+ctx.Tr("certificate")+".pdf")
	ctx.Resp.Write(file.Bytes())
}

func GenerateFont() error {
	fontDir := "static/fonts/"
	return gofpdf.MakeFont(fontDir+"zillaslab.ttf", fontDir+"cp1252.map", fontDir, nil, true)
}

func generateCertificate(file io.Writer, user storage.User, ctx *macaron.Context) error {
	fontDir := "static/fonts/"

	pdf := gofpdf.New("P", "mm", "A4", fontDir)
	utf8tr := pdf.UnicodeTranslatorFromDescriptor("cp1252")
	pdf.AddFont("Zilla Slab", "", "zillaslab.json")
	pdf.AddPage()

	pdf.Ln(39)
	pdf.SetFont("Zilla Slab", "", 109)
	pdf.SetTextColor(196, 187, 69)
	pdf.CellFormat(0, 32, utf8tr(ctx.Tr("certificate")), "", 1, "C", false, 0, "")

	pdf.Ln(39)
	pdf.SetFont("Zilla Slab", "", 35)
	pdf.SetTextColor(34, 34, 34)
	pdf.CellFormat(0, 10, utf8tr(ctx.Tr("awarded_to")), "", 1, "C", false, 0, "")

	pdf.Ln(25)
	pdf.SetFont("Zilla Slab", "", 38)
	pdf.SetTextColor(196, 187, 69)
	pdf.CellFormat(0, 10, utf8tr(user.Email), "", 1, "C", false, 0, "")

	scoreStr := ctx.Tr("certificate_score_pre") + " " + strconv.Itoa(user.Score) + " " + ctx.Tr("certificate_score_post")
	pdf.Ln(25)
	pdf.SetFont("Zilla Slab", "", 35)
	pdf.SetTextColor(34, 34, 34)
	pdf.CellFormat(0, 10, utf8tr(scoreStr), "", 1, "C", false, 0, "")

	pdf.Ln(50)
	pdf.SetFont("Zilla Slab", "", 27)
	pdf.CellFormat(0, 10, utf8tr(ctx.Tr("service")), "", 1, "C", false, 0, "")

	err := pdf.Output(file)
	pdf.Close()

	return err
}

func certificateReady() bool {
	_, month, day := time.Now().Date()
	return (month == time.December && day > 25) || (int(month) < config.Config.Server.ResetMonth)
}
