package routes

import macaron "gopkg.in/macaron.v1"

func Home(ctx *macaron.Context) {
	ctx.HTML(200, "home")
}
