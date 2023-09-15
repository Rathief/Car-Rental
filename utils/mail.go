package utils

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "rafie.fadlurahman@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	pass, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		pass,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
