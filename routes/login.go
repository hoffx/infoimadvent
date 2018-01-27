package routes

import (
	"log"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

func Login(ctx *macaron.Context, log *log.Logger, uStorer *storage.UserStorer, sess session.Store) {
	defer ctx.HTML(200, "login")
	defer parseURL(ctx)

	type Config struct {
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}
	ctx.Data["IsAdvent"] = isAdvent()

	if ctx.Req.Method == "GET" {
		value := sess.Get("user")
		sUser, ok := value.(storage.User)
		if ok && sUser.Active {
			ctx.Data["LoggedIn"] = true
		}
	} else {
		fEmail := ctx.Req.FormValue("email")
		fPw := ctx.Req.FormValue("pw")

		ctx.Data["Email"] = fEmail

		user, err := uStorer.Get(map[string]interface{}{"email": fEmail})
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrDB)
			log.Println(err)
			return
		} else if user.Email == "" {
			ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials)
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(fPw))
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials)
			return
		}
		if !user.Confirmed {
			ctx.Data["Error"] = ctx.Tr(ErrNotConfirmed)
			return
		}
		user.Active = true

		err = sess.Set("user", user)
		if err != nil {
			user.Active = false
			ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
			log.Println(err)
			return
		}

		err = uStorer.Put(user)
		if err != nil {
			user.Active = false
			ctx.Data["Error"] = ctx.Tr(ErrDB)
			log.Println(err)
			return
		}

		ctx.Data["LoggedIn"] = true
		ctx.Data["Message"] = ctx.Tr(MessLoggedIn)
	}

}

func parseURL(ctx *macaron.Context) {
	query := ctx.Req.URL.Query()
	for k, v := range query {
		if len(v) == 1 {
			ctx.Data[k] = v[0]
		} else if len(v) > 1 {
			ctx.Data[k] = v
		}
	}
}
