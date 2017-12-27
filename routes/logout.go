package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Logout(ctx *macaron.Context, sess session.Store) {
	sess.Set("user", storage.User{})
	ctx.Redirect("/login", 200)
}
