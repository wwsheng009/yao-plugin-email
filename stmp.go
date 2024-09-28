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

	tousers := make([]string, len(message.To))
	for i, recipient := range message.To {
		tousers[i] = m.FormatAddress(recipient.Address, recipient.Name)
	}

	m.SetHeader("To", tousers...)
	if len(message.CC) > 0 {
		ccusers := make([]string, len(message.CC))
		for i, recipient := range message.CC {
			ccusers[i] = m.FormatAddress(recipient.Address, recipient.Name)
		}
		m.SetHeader("Cc", ccusers...)
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
