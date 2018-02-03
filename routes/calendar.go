package routes

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"
)

type templDay struct {
	Link    string
	Opened  bool
	Date    int
	Current bool
	Locked  bool
}

func Calendar(ctx *macaron.Context, log *log.Logger, sess session.Store) {

	value := sess.Get("user")
	// protected therefore user must exist
	user := value.(storage.User)

	tDays := make([]templDay, 0)
	for i, d := range user.Days {
		_, month, day := time.Now().Date()

		var opened, current, locked bool

		// TODO: change back to december after testing
		if d != storage.None || (i+1 < day && month == time.February) || month != time.February {
			opened = true
		}
		// TODO: change back to december after testing
		if month == time.February && day == i+1 {
			current = true
		}
		// TODO: change back to december after testing
		if month != time.February || day < i+1 {
			locked = true
		}
		tDays = append(tDays, templDay{"/day/" + strconv.Itoa(i+1), opened, i + 1, current, locked})
	}
	randomize(&tDays)
	ctx.Data["Days"] = tDays
	ctx.HTML(200, "calendar")
}

func randomize(days *[]templDay) {
	t := *days
	l := len(t)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, d := range t {
		j := r.Intn(l - 1)
		t[i] = t[j]
		t[j] = d
	}
	days = &t
}
