package utils

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "piyushrajak143@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "piyushrajak143@gmail.com", "jhthfkxkalynzvpi") // no spaces in app password

	err := d.DialAndSend(m)
	if err != nil {
		log.Fatal("Failed to send email:", err)
	}
	return err
}
