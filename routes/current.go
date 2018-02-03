package routes

import (
	"log"
	"strconv"
	"time"

	macaron "gopkg.in/macaron.v1"
)

func Current(ctx *macaron.Context, log *log.Logger) {
	_, m, d := time.Now().Date()
	// TODO: change back to december after testing
	if m == time.February && d > 0 && d <= 24 {
		ctx.Redirect("/day/" + strconv.Itoa(d) + "?" + ctx.Req.URL.RawQuery)
	} else {
		ctx.Error(406, ctx.Tr(ErrIllegalDate))
		ctx.Redirect("/", 406)
	}
}
