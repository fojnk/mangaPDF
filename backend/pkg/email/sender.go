package email

import (
	"os"

	"github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

func SendHtmlEmail(to []string, subject string, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("FROM_EMAIL"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(viper.GetString("smtp.server"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	return d.DialAndSend(m)
}
