package routes

import (
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

func PublicReady(ctx *macaron.Context, qStorer *storage.QuestStorer) {
	if !qStorer.Complete {
		ctx.Error(503, ctx.Tr(ErrNotReady))
	}
}
