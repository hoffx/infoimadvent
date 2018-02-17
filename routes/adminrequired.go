package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

// RequireAdmin ensures the routes "/upload" and "/overview" can
// only be accessed by the admin
func RequireAdmin(ctx *macaron.Context, sess session.Store) {
	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if !ok || !sUser.Active || !sUser.IsAdmin {
		ctx.Redirect("/", 401)
	}
}
