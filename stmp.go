package main

import (
	"crypto/tls"
	"fmt"

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
	// 验证账户信息
	if account.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if account.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if account.Server == "" {
		return fmt.Errorf("SMTP server is required")
	}
	if account.Port <= 0 {
		return fmt.Errorf("SMTP port must be greater than 0")
	}

	// 验证消息内容
	if len(message.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	if message.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if message.Body == "" {
		return fmt.Errorf("email body is required")
	}

	// 验证收件人地址
	for i, recipient := range message.To {
		if recipient.Address == "" {
			return fmt.Errorf("recipient %d: email address is required", i+1)
		}
	}

	// 验证抄送地址
	for i, recipient := range message.CC {
		if recipient.Address == "" {
			return fmt.Errorf("CC recipient %d: email address is required", i+1)
		}
	}

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
