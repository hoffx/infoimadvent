package routes

import (
	"log"

	macaron "gopkg.in/macaron.v1"
)

func Register(ctx *macaron.Context, log *log.Logger) {

	type Config struct {
		MinYL, MaxYL int
	}

	ctx.Data["Config"] = Config{1, 12}

	ctx.HTML(200, "register")
}
