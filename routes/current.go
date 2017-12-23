package routes

import (
	"log"
	"strconv"
	"time"

	macaron "gopkg.in/macaron.v1"
)

func Current(ctx *macaron.Context, log *log.Logger) {
	_, m, d := time.Now().Date()
	if m == time.December && d >= 0 && d <= 24 {
		ctx.Redirect("/day/" + strconv.Itoa(d))
	} else {
		ctx.Redirect("/")
	}
}
