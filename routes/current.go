package routes

import (
	"log"
	"strconv"
	"time"

	"github.com/hoffx/infoimadvent/config"
	macaron "gopkg.in/macaron.v1"
)

// Current handles the route "/day". It is a shortcut for "/day/[current_day]"
func Current(ctx *macaron.Context, log *log.Logger) {
	_, m, d := time.Now().Date()
	if m == config.Config.Server.Advent && d > 0 && d <= 24 {
		ctx.Redirect("/day/" + strconv.Itoa(d) + "?" + ctx.Req.URL.RawQuery)
	} else {
		ctx.Error(406, ctx.Tr(ErrIllegalDate))
		ctx.Redirect("/", 406)
	}
}
