package services

import (
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	gomail "gopkg.in/gomail.v2"
	macaron "gopkg.in/macaron.v1"
)

func (s *DBStorage) sendRewardMail(ctx *macaron.Context, user storage.User) (err error) {
	ctx.Data["User"] = user

	r := macaron.TplRender{}

	mailBody, err := r.HTMLString("rewardmail", ctx.Data)
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
