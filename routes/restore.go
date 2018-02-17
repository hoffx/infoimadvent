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

// Restore handles the route "/restore"
func Restore(ctx *macaron.Context, log *log.Logger, uStorer *storage.UserStorer) {
	email := ctx.Req.FormValue("email")

	// try to find user

	user, err := uStorer.Get(map[string]interface{}{"email": email})
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrDB))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}
	if user.Email == "" {
		// user not found, but stating an e-mail was sent, so that this information can't abused for bruteforce-hacking
		ctx.Redirect("/login?Email="+user.Email+"&Message="+ctx.Tr(MessRestoreMailSent), 302)
		return
	}
	if !user.Confirmed {
		ctx.Redirect("/login?Error="+ErrNotConfirmed, 302)
		return
	}

	pw, err := gostrgen.RandGen(int(config.Config.Grades.Max), gostrgen.LowerUpperDigit, "", "")
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}

	user.Hash = string(hash)
	err = uStorer.Put(user)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrDB))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}

	// send restore email

	ctx.Data["User"] = user
	ctx.Data["Password"] = pw

	mailBody, err := ctx.HTMLString("restoremail", ctx.Data)
	if err != nil {
		ctx.Error(500, ctx.Tr(ErrUnexpected))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Mail.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", ctx.Tr("restore_mail_header"))
	m.SetBody("text/html", mailBody)

	d := gomail.NewDialer(config.Config.Mail.Address, config.Config.Mail.Port, config.Config.Mail.Username, config.Config.Mail.Password)

	if err := d.DialAndSend(m); err != nil {
		ctx.Error(500, ctx.Tr(ErrMail))
		log.Println(err)
		ctx.Redirect("/login", 500)
		return
	}

	ctx.Redirect("/login?Email="+user.Email+"&Message="+ctx.Tr(MessRestoreMailSent), 302)
}
