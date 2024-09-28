package main

import (
	"net/mail"
	"time"
)

type Account struct {
	Server   string `json:"server"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}
type Recipient struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
}
type Message struct {
	From        string      `json:"from"`
	To          []Recipient `json:"to"`
	CC          []Recipient `json:"cc,omitempty"`
	Subject     string      `json:"subject"`
	Body        string      `json:"body"`
	Attachments []string    `json:"attachments,omitempty"`
}

type Email struct {
	Account  Account   `json:"account"`
	Messages []Message `json:"messages"`
}

type BodyPart struct {
	ContentType      string            `json:"content_type"`
	ContentTypeValue map[string]string `json:"content_type_value"`
	Disposition      string            `json:"disposition"`
	DispositionValue map[string]string `json:"disposition_value"`
	Charset          string            `json:"charset"`
	FileName         string            `json:"file_name"`
	FileSize         int64             `json:"file_size"`
	SavedFileName    string            `json:"saved_file_name"`
	SavedFilePath    string            `json:"saved_file_path"`
	Encoding         string            `json:"encoding"`
	ContentId        string            `json:"content_id"`
	Centent          string            `json:"centent"`
	Attachment       string            `json:"attachment"`
}
type MessageReceived struct {
	From        []mail.Address `json:"from"`
	To          []mail.Address `json:"to"`
	Subject     string         `json:"subject"`
	Body        []BodyPart     `json:"body"`
	Attachments []string       `json:"attachments,omitempty"`
	Date        time.Time      `json:"date"`
	Error       string         `json:"error"`
	MessageId   string         `json:"message_id"`
	Folder      string         `json:"folder"`
	Uid         uint32         `json:"uid"`
}
