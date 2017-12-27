package routes

import (
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Confirm(ctx *macaron.Context, storer *storage.Storer) {
	query := ctx.Req.URL.Query()
	// check if request is formatted correctly
	if query["user"] != nil && query["token"] != nil {
		// get according user
		user, err := storer.Get(query["user"][0])
		if err != nil {
			ctx.Error(500, ErrDB.Error())
			ctx.Redirect("/", 500)
			return
		}
		// check if token and user status are ok
		if !user.Confirmed && user.ConfirmationToken == query["token"][0] {
			user.ConfirmationToken = ""
			user.Confirmed = true
			err = storer.Put(user)
			if err != nil {
				ctx.Error(500, ErrDB.Error())
				ctx.Redirect("/", 500)
				return
			}
		} else {
			// redirect user (that was messing around with the links) to home-page
			ctx.Error(406, ErrWrongCredentials.Error())
			ctx.Redirect("/", 406)
			return
		}
	} else {
		ctx.Error(406, ErrWrongCredentials.Error())
		ctx.Redirect("/", 406)
		return
	}
	// TODO: automated login
	ctx.Redirect("/login", 302)
}
