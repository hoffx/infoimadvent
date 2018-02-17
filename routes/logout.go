package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

// Logout handles the route "/logout"
func Logout(ctx *macaron.Context, uStorer *storage.UserStorer, sess session.Store) {
	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if ok {
		sUser.Active = false
		err := sess.Set("user", storage.User{})
		if err != nil {
			ctx.Error(500, ctx.Tr(ErrUnexpected))
			ctx.Redirect("/login", 500)
			return
		}
		err = uStorer.Put(sUser)
		if err != nil {
			ctx.Error(500, ctx.Tr(ErrDB))
			ctx.Redirect("/login", 500)
			return
		}
	}
	ctx.Redirect("/login?Message="+ctx.Tr(MessLoggedOut), 302)
}
