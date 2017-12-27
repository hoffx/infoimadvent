package routes

import (
	"errors"
	"log"
	"strconv"

	"github.com/elgs/gostrgen"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	macaron "gopkg.in/macaron.v1"

	"gopkg.in/gomail.v2"
)

var ErrUnequalPasswords = errors.New("passwords_unequal")
var ErrUserExists = errors.New("user_exists")
var ErrMail = errors.New("mail_error")
var ErrWrongGrade = errors.New("wrong_grade_error")

func Register(ctx *macaron.Context, log *log.Logger, storer *storage.Storer) {
	defer ctx.HTML(200, "register")

	// TODO: handle form refill on failure

	type Config struct {
		MinYL, MaxYL   uint
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Grades.Min, config.Config.Grades.Max, config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}

	if ctx.Req.Method == "GET" {
		// handle get requests

		// TODO: check if user is logged in using session and print logout link instead of register form
		if false {
			ctx.Data["LoggedIn"] = true
		}
	} else {
		// handle post requests

		// check input
		fEmail := ctx.Req.FormValue("email")
		fPw := ctx.Req.FormValue("pw1")
		fPw2 := ctx.Req.FormValue("pw2")
		g, err := strconv.Atoi(ctx.Req.FormValue("grade"))
		if err != nil {
			ctx.Data["Error"] = err.Error()
			return
		}
		fGrade := uint(g)
		if fGrade < config.Config.Grades.Min || fGrade > config.Config.Grades.Max {
			ctx.Data["Error"] = ErrWrongGrade.Error()
			return
		}
		if fPw != fPw2 {
			ctx.Data["Error"] = ErrUnequalPasswords.Error()
			return
		}

		user, err := storer.Get(fEmail)
		if err != nil {
			ctx.Data["Error"] = ErrDB.Error()
			return
		}
		if user.Email != "" {
			ctx.Data["Error"] = ErrUserExists.Error()
		} else {
			// request accepted -> generating new user

			// create user and write to db
			confirmationToken, err := gostrgen.RandGen(40, gostrgen.LowerUpperDigit, "", "")
			if err != nil {
				ctx.Data["Error"] = ErrUnexpected.Error()
				return
			}
			user = storage.User{fEmail, fPw, uint(fGrade), false, false, confirmationToken}
			err = storer.Create(user)
			if err != nil {
				ctx.Data["Error"] = ErrDB.Error()
				return
			}

			// send confirmation email
			m := gomail.NewMessage()
			m.SetHeader("From", config.Config.Mail.Sender)
			m.SetHeader("To", user.Email)
			m.SetHeader("Subject", ctx.Tr("confirmation_mail_header"))
			m.SetBody("text/html", ctx.Tr("confirmation_mail_body")+`<a href="http://`+config.Config.Server.Address+`/confirm?user=`+user.Email+`&token=`+confirmationToken+`" > http://`+config.Config.Server.Address+`/confirm?user=`+user.Email+`&token=`+confirmationToken+`</a>`)

			d := gomail.NewDialer(config.Config.Mail.Address, config.Config.Mail.Port, config.Config.Mail.Username, config.Config.Mail.Password)

			if err := d.DialAndSend(m); err != nil {
				ctx.Data["Error"] = ErrMail.Error()
			}
		}
	}
}
