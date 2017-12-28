package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Protect(ctx *macaron.Context, sess session.Store) {
	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if !ok || (ok && !sUser.Active) {
		ctx.Redirect("/", 401)
	}
}
