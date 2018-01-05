package routes

import (
	"log"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

func Account(ctx *macaron.Context, log *log.Logger, sess session.Store, uStorer *storage.UserStorer) {
	defer ctx.HTML(200, "account")

	type Config struct {
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}

	if ctx.Req.Method == "GET" {
		if mode := ctx.Req.URL.Query()["mode"]; len(mode) > 0 && mode[0] == "changepw" {
			ctx.Data["ChangePw"] = true
		} else {
			user := sess.Get("user").(storage.User)
			user, err := uStorer.Get(map[string]interface{}{"email": user.Email})
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrDB)
				log.Println(err)
				return
			}

			ctx.Data["Score"] = user.Score

			err = sess.Set("user", user)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
				log.Println(err)
				return
			}
		}
	} else {
		ctx.Data["ChangePw"] = true

		user := sess.Get("user").(storage.User)
		user, err := uStorer.Get(map[string]interface{}{"email": user.Email})
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrDB)
			log.Println(err)
			return
		}
		ctx.Data["Score"] = user.Score

		fPwOld := ctx.Req.FormValue("pwOld")
		fPw := ctx.Req.FormValue("pw1")
		fPw2 := ctx.Req.FormValue("pw2")

		ctx.Data["PwOld"] = fPwOld

		if fPw != fPw2 {
			ctx.Data["Error"] = ctx.Tr(ErrUnequalPasswords)
			return
		}

		ctx.Data["Pw"] = fPw

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fPwOld))
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials)
			return
		}

		pwBytes, err := bcrypt.GenerateFromPassword([]byte(fPw), bcrypt.DefaultCost)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
			return
		}
		user.Password = string(pwBytes)

		err = uStorer.Put(user)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrDB)
			log.Println(err)
			return
		}

		err = sess.Set("user", user)
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
			log.Println(err)
			return
		}

		ctx.Data["Message"] = ctx.Tr(MessChangedPassword)
	}
}
