package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func Logout(ctx *macaron.Context) {
	// TODO: implement logout
	ctx.Redirect("/login", 200)
}
