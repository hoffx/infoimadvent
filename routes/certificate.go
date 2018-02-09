package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Certificate(ctx *macaron.Context, sess session.Store) {
	user := sess.Get("user").(storage.User)

	ctx.Data["Email"] = user.Email
	ctx.Data["Score"] = user.Score

	// htmlBody, err := ctx.HTMLString("certificate", ctx.Data)

	// if err != nil {
	// 	ctx.Error(500, ctx.Tr(ErrUnexpected))
	// 	log.Println(err)
	// 	ctx.Redirect("/certificate", 500)
	// 	return
	// }

	ctx.HTML(200, "certificate")
}
