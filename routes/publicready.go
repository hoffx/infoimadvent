package routes

import (
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func PublicReady(ctx *macaron.Context, dStorer *storage.DocumentStorer) {
	if !dStorer.Complete {
		ctx.Error(503, ctx.Tr(ErrNotReady))
	}
}
