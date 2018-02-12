package routes

import (
	"time"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Home(ctx *macaron.Context, sess session.Store) {
	query := ctx.Req.URL.Query()

	if query["cookies"] != nil && query["cookies"][0] == "accepted" {
		ctx.SetCookie("cookies", "accepted", 1<<31-1)
		ctx.Data["Cookies"] = true
	} else {
		ctx.Data["Cookies"] = ctx.GetCookie("cookies") == "accepted"
	}

	value := sess.Get("user")
	sUser, ok := value.(storage.User)
	if ok && sUser.Active {
		ctx.Data["LoggedIn"] = true
	}
	ctx.Data["IsAdvent"] = isAdvent()
	ctx.HTML(200, "home")
}

func isAdvent() bool {
	_, m, d := time.Now().Date()
	if m == time.February && d >= 1 && d <= 24 {
		return true
	}
	return false
}
