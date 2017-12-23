package routes

import (
	"time"

	macaron "gopkg.in/macaron.v1"
)

func Day(ctx *macaron.Context) {
	num := ctx.ParamsInt("day")
	if num < 1 || num > 24 {
		ctx.HTML(404, "404")
		return
	}
	_, m, d := time.Now().Date()
	if m != time.December || d < num {
		ctx.HTML(403, "wait")
		return
	}
	ctx.HTML(200, "day")
}
