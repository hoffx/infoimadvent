package routes

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/go-macaron/session"
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

func generateCertificate(file io.Writer, user storage.User, ctx *macaron.Context) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, ctx.Tr("score")+": "+strconv.Itoa(user.Score))

	err := pdf.Output(file)
	pdf.Close()

	return err
}

func certificateReady() bool {
	// TODO: change back after testing
	//_, month, day := time.Now().Date()

	//return month != time.February && day > 24
	return true
}
