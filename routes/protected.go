package routes

import (
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Protect(ctx *macaron.Context, user *storage.User) {
	if !user.Active {
		//sess.Set("redirect", ctx.Req.RequestURI)
		ctx.Redirect("/", 302)
	}
}
