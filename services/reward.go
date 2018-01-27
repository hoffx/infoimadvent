package services

import (
	"github.com/go-macaron/i18n"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	gomail "gopkg.in/gomail.v2"
	macaron "gopkg.in/macaron.v1"
)

func SendRewardMail(ctx *macaron.Context, user storage.User) (err error) {
	//var data map[string]interface{}
	//data["User"] = user

	i18n.I18n(i18n.Options{
		Directory: "locales",
		Langs:     []string{"de-DE", "en-US"},
		Names:     []string{"Deutsch", "Englisch"},
	})

	//r := macaron.TplRender{}

	mailBody, err := ctx.HTMLString("rewardmail", ctx.Data)
	if err != nil {
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Mail.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", ctx.Tr("reward_mail_header"))
	m.SetBody("text/html", mailBody)

	d := gomail.NewDialer(config.Config.Mail.Address, config.Config.Mail.Port, config.Config.Mail.Username, config.Config.Mail.Password)

	if err = d.DialAndSend(m); err != nil {
		return
	}
	return
}
