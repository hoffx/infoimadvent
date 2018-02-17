package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

// Protect ensures the routes "/logout", "/account", "/certificate",
// "/calendar" and "/day" can only be accessed when logged in
func Protect(ctx *macaron.Context, sess session.Store) {
	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if !ok || !sUser.Active {
		ctx.Redirect("/", 401)
	}
}
