package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

var MessLoggedIn = "login_success"

func Confirm(ctx *macaron.Context, storer *storage.Storer, sess session.Store) {
	query := ctx.Req.URL.Query()
	var user storage.User
	var err error
	// check if request is formatted correctly
	if query["user"] != nil && query["token"] != nil {
		// get according user
		user, err = storer.Get(query["user"][0])
		if err != nil {
			ctx.Error(500, ErrDB.Error())
			ctx.Redirect("/", 500)
			return
		}
		// check if token and user status are ok
		if user.Email != "" && !user.Confirmed && user.ConfirmationToken == query["token"][0] {
			user.ConfirmationToken = ""
			user.Confirmed = true
			user.Active = true

			sess.Set("user", user)

			err = storer.Put(user)
			if err != nil {
				ctx.Error(500, ErrDB.Error())
				ctx.Redirect("/", 500)
				return
			}
		} else {
			// redirect user (that was messing around with the link) to home-page
			ctx.Error(406, ErrWrongCredentials.Error())
			ctx.Redirect("/", 406)
			return
		}
	} else {
		ctx.Error(406, ErrWrongCredentials.Error())
		ctx.Redirect("/", 406)
		return
	}
	ctx.Redirect("/login?Message="+MessLoggedIn, 302)
}
