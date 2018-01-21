package routes

import (
	"encoding/base64"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/elgs/gostrgen"
	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"

	"gopkg.in/gomail.v2"
)

func Register(ctx *macaron.Context, log *log.Logger, uStorer *storage.UserStorer, sess session.Store) {
	defer ctx.HTML(200, "register")

	// TODO: handle form refill on failure

	type Config struct {
		MinYL, MaxYL   uint
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Grades.Min, config.Config.Grades.Max, config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}
	ctx.Data["Teacher"] = false

	if ctx.Req.Method == "GET" {
		// handle get requests

		value := sess.Get("user")
		sUser, ok := value.(storage.User)
		if ok && sUser.Active {
			ctx.Data["LoggedIn"] = true
		}
	} else {
		// handle post requests

		// check input
		fEmail := ctx.Req.FormValue("email")
		fPw := ctx.Req.FormValue("pw1")
		fPw2 := ctx.Req.FormValue("pw2")

		ctx.Data["Email"] = fEmail
		ctx.Data["Pw"] = fPw

		t := ctx.Req.FormValue("teacher") == "on"

		ctx.Data["Teacher"] = t

		g, err := strconv.Atoi(ctx.Req.FormValue("grade"))

		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
			return
		}

		ctx.Data["Grade"] = g

		fGrade := uint(g)
		if fGrade < config.Config.Grades.Min || fGrade > config.Config.Grades.Max {
			ctx.Data["Error"] = ctx.Tr(ErrWrongGrade)
			ctx.Data["Grade"] = nil
			return
		}
		if fPw != fPw2 {
			ctx.Data["Error"] = ctx.Tr(ErrUnequalPasswords)
			ctx.Data["Pw"] = nil
			return
		}

		user, err := uStorer.Get(map[string]interface{}{"email": fEmail})
		if err != nil {
			ctx.Data["Error"] = ctx.Tr(ErrDB)
			log.Println(err)
			return
		}
		if user.Email != "" {
			ctx.Data["Error"] = ctx.Tr(ErrUserExists)
			ctx.Data["Grade"] = nil
		} else {
			// request accepted -> generating new user

			// create user and write to db
			confirmationToken, err := gostrgen.RandGen(40, gostrgen.LowerUpperDigit, "", "")
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
				log.Println(err)
				return
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(fPw), bcrypt.DefaultCost)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
				log.Println(err)
				return
			}

			user = storage.User{fEmail, string(hash), uint(fGrade), false, false, confirmationToken, t, make([]int, 24), 0, false}
			err = uStorer.Create(user)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrDB)
				log.Println(err)
				return
			}

			// send confirmation email
			m := gomail.NewMessage()
			m.SetHeader("From", config.Config.Mail.Sender)
			m.SetHeader("To", user.Email)
			m.SetHeader("Subject", ctx.Tr("confirmation_mail_header"))
			encodedMail := base64.StdEncoding.EncodeToString([]byte(user.Email))
			m.SetBody("text/html", ctx.Tr("confirmation_mail_body")+`<a href="http://`+config.Config.Server.Address+`/confirm?user=`+encodedMail+`&token=`+confirmationToken+`" > http://`+config.Config.Server.Address+`/confirm?user=`+encodedMail+`&token=`+confirmationToken+`</a>`)

			d := gomail.NewDialer(config.Config.Mail.Address, config.Config.Mail.Port, config.Config.Mail.Username, config.Config.Mail.Password)

			if err := d.DialAndSend(m); err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrMail)
				log.Println(err)

				err = uStorer.Delete(map[string]interface{}{"email": user.Email})
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				}

				return
			}

			ctx.Data["Message"] = ctx.Tr(MessConfirmMailSent)
		}
	}
}
