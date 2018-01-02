package routes

import (
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func Account(ctx *macaron.Context, sess session.Store, uStorer *storage.UserStorer) {
	if ctx.Req.Method == "GET" {

	} else {

	}
}
