package main

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

func sendStmpEmails(email Email) []error {
	errs := make([]error, 0)
	for _, v := range email.Messages {
		err := sendStmpEmail(email.Account, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func sendStmpEmail(account Account, message Message) error {
	m := gomail.NewMessage()
	m.SetHeader("From", message.From)
	if message.From == "" {
		m.SetHeader("From", account.Username)
	}
	m.SetHeader("To", message.To)
	if len(message.CC) > 0 {
		m.SetHeader("Cc", message.CC...)
	}
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", message.Body)

	for _, attachment := range message.Attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(account.Server, account.Port, account.Username, account.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
