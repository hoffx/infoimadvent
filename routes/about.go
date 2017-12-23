package routes

import (
	"log"

	macaron "gopkg.in/macaron.v1"
)

func About(ctx *macaron.Context, log *log.Logger) {
	ctx.Redirect("/")
}
