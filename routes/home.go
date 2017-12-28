package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Home(ctx *macaron.Context, sess session.Store) {
	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if ok && sUser.Active {
		ctx.Data["LoggedIn"] = true
	}
	ctx.HTML(200, "home")
}
