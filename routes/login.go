package routes

import (
	"errors"
	"log"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

var ErrDB = errors.New("database_error")
var ErrUnexpected = errors.New("unexpected_error")
var ErrWrongCredentials = errors.New("wrong_credentials_error")
var ErrNotConfirmed = errors.New("user_not_confirmed")

func Login(ctx *macaron.Context, log *log.Logger, storer *storage.Storer, sess session.Store) {
	defer ctx.HTML(200, "login")
	defer parseURL(ctx)

	type Config struct {
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}

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

		user, err := storer.Get(fEmail)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrDB.Error())
			log.Println(err)
			return
		} else if user.Email == "" {
			ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials.Error())
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fPw))
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials.Error())
			return
		}
		if !user.Confirmed {
			ctx.Data["Error"] = ctx.Tr(ErrNotConfirmed.Error())
			return
		}
		ctx.Data["LoggedIn"] = true
		ctx.Data["Message"] = ctx.Tr(MessLoggedIn)
		user.Active = true

		sess.Set("user", user)
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
