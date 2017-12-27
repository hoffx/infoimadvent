package routes

import (
	"errors"

	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

var ErrDB = errors.New("database_error")
var ErrUnexpected = errors.New("unexpected_error")
var ErrWrongCredentials = errors.New("wrong_credentials_error")
var ErrNotConfirmed = errors.New("user_not_confirmed")

func Login(ctx *macaron.Context, storer *storage.Storer) {
	defer ctx.HTML(200, "login")

	type Config struct {
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}

	if ctx.Req.Method == "GET" {
		// TODO: check if user is logged in using session and print logout link instead of login form
		if false {
			ctx.Data["LoggedIn"] = true
		}
	} else {
		fEmail := ctx.Req.FormValue("email")
		fPw := ctx.Req.FormValue("pw")
		user, err := storer.Get(fEmail)
		if err != nil {
			ctx.Data["Error"] = ErrDB.Error()
			return
		} else if user.Email == "" {
			ctx.Data["Error"] = ErrWrongCredentials.Error()
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fPw))
		if err != nil {
			ctx.Data["Error"] = ErrWrongCredentials.Error()
			return
		}
		if !user.Confirmed {
			ctx.Data["Error"] = ErrNotConfirmed.Error()
			return
		}
		ctx.Data["LoggedIn"] = true
		user.Active = true

		// TODO: create session
	}

}
