package routes

import (
	"log"

	macaron "gopkg.in/macaron.v1"
)

func Calendar(ctx *macaron.Context, log *log.Logger) {
	ctx.Redirect("/")
}
