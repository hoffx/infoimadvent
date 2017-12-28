package routes

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/elgs/gostrgen"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	gomail "gopkg.in/gomail.v2"
	macaron "gopkg.in/macaron.v1"
)

var MessRestoreMailSent = "restore_mail_sent"

func Restore(ctx *macaron.Context, log *log.Logger, storer *storage.Storer) {
	email := ctx.Req.FormValue("email")

	// try to find user

	user, err := storer.Get(email)
	if err != nil {
		ctx.Error(500, ErrDB.Error())
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	} else if user.Email == "" {
		// user not found, but stating an e-mail was sent, so that this information can't abused for bruteforce-hacking
		ctx.Redirect("/login?Email="+user.Email+"&Message="+ctx.Tr(MessRestoreMailSent), 302)
		return
	} else if !user.Confirmed {
		ctx.Redirect("/login?Error="+ErrNotConfirmed.Error(), 302)
		return
	} else {
		pw, err := gostrgen.RandGen(int(config.Config.Grades.Max), gostrgen.LowerUpperDigit, "", "")
		if err != nil {
			ctx.Error(500, ErrUnexpected.Error())
			log.Println(err)
			ctx.Redirect("/login", 500)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			ctx.Error(500, ErrUnexpected.Error())
			log.Println(err)
			ctx.Redirect("/login", 500)
			return
		}

		user.Password = string(hash)
		err = storer.Put(user)
		if err != nil {
			ctx.Error(500, ErrDB.Error())
			log.Println(err)
			ctx.Redirect("/login", 500)
			return
		}

		// TODO: write restore mail's text

		m := gomail.NewMessage()
		m.SetHeader("From", config.Config.Mail.Sender)
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", ctx.Tr("restore_mail_header"))
		m.SetBody("text/html", ctx.Tr("restore_mail_body")+pw)

		d := gomail.NewDialer(config.Config.Mail.Address, config.Config.Mail.Port, config.Config.Mail.Username, config.Config.Mail.Password)

		if err := d.DialAndSend(m); err != nil {
			ctx.Error(500, ErrMail.Error())
			log.Println(err)
			ctx.Redirect("/login", 500)
			return
		}
	}

	ctx.Redirect("/login?Email="+user.Email+"&Message="+ctx.Tr(MessRestoreMailSent), 302)
}
