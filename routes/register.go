package routes

import (
	"log"

	macaron "gopkg.in/macaron.v1"
)

func Register(ctx *macaron.Context, log *log.Logger) {
	ctx.HTML(200, "register")
}
