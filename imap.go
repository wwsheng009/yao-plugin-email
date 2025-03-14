package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func getUniqueFolder(messageId string) string {

	// Create an MD5 hash of the string
	hash := md5.Sum([]byte(messageId))

	// Convert the hash to a hexadecimal string
	hashStr := hex.EncodeToString(hash[:])

	// Use the hash as the unique part of the file name
	// fileName := fmt.Sprintf("%s.txt", hashStr)
	return hashStr
}
func makeFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			// fmt.Println("Failed to create folder:", err)
			return err
		}
	}
	return nil
}
func convertUtf8String(encodedString string) string {
	encodedString = strings.TrimPrefix(encodedString, "=?utf-8?B?")
	encodedString = strings.TrimSuffix(encodedString, "=?=")
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return encodedString
	}
	utf8String := string(decodedBytes)
	return utf8String
}
func connectToServer(username, password, server string, port int) (*client.Client, error) {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", server, port), nil)
	if err != nil {
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	charset.RegisterEncoding("gb2312", simplifiedchinese.GB18030)

	return c, nil
}

func fetchImapEmails(imapClient *client.Client, email Email) ([]MessageReceived, error) {

	emailObjects := make([]MessageReceived, 0)

	// Select the mbox you want to read
	mbox, err := imapClient.Select("INBOX", false)
	if err != nil {
		return emailObjects, err
	}

	// Get the last message
	if mbox.Messages == 0 {
		return emailObjects, err
	}
	// criteria to search for unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	uids, err := imapClient.UidSearch(criteria)
	if err != nil {
		return emailObjects, err
	}
	if len(uids) == 0 {
		return emailObjects, err
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	// used to collect the seqnum
	seqSet2 := new(imap.SeqSet)

	// seqSet.AddRange(1, mbox.Messages)
	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- imapClient.UidFetch(seqSet, items, messages)
	}()

	for msg := range messages {
		emailObj := MessageReceived{}

		if msg != nil {
			seqSet2.AddNum(msg.SeqNum)
			emailObj.Uid = msg.Uid
		} else {
			emailObj.Error = "Server didn't returned message"
			emailObjects = append(emailObjects, emailObj)
			continue
		}

		r := msg.GetBody(&section)
		if r == nil {
			emailObj.Error = "Server didn't returned message body"
			emailObjects = append(emailObjects, emailObj)
			continue
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			emailObj.Error = err.Error()
			emailObjects = append(emailObjects, emailObj)
			continue
		}

		// Print some info about the message
		header := mr.Header

		if messageId, err := header.MessageID(); err == nil {
			emailObj.MessageId = messageId
		}
		folder := getUniqueFolder(emailObj.MessageId)
		if email.Folder == "" {
			email.Folder = "data/upload/emails"
		}
		folder = email.Folder + "/" + folder
		emailObj.Folder = folder

		if date, err := header.Date(); err == nil {
			emailObj.Date = date
		}
		if from, err := header.AddressList("From"); err == nil {
			for _, line := range from {
				if strings.Contains(line.Name, "utf-8") {
					line.Name = convertUtf8String(line.Name)
				}
				emailObj.From = append(emailObj.From, *line)
			}
		}
		if to, err := header.AddressList("To"); err == nil {
			// log.Println("To:", to)
			for _, line := range to {
				if strings.Contains(line.Name, "utf-8") {
					line.Name = convertUtf8String(line.Name)
				}
				emailObj.To = append(emailObj.To, *line)
			}
		}
		if subject, err := header.Subject(); err == nil {
			if strings.Contains(subject, "utf-8") {
				subject = convertUtf8String(subject)
			}
			emailObj.Subject = subject
		}

		headers := make([]BodyPart, 0)
		for {
			headerPart := BodyPart{}

			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				// log.Fatal(err)
				emailObj.Error = err.Error()
				emailObjects = append(emailObjects, emailObj)
				break
			}
			headerPart.Encoding = p.Header.Get("Content-Transfer-Encoding")
			headerPart.ContentId = p.Header.Get("Content-Id")

			switch h := p.Header.(type) {
			case *mail.InlineHeader:

				// This is the message's text (can be plain-text or HTML)
				contentType, params, err := h.ContentType()
				if err == nil {
					headerPart.ContentType = contentType
					headerPart.ContentTypeValue = params
				}
				disp, params, err := h.ContentDisposition()
				if err == nil {
					headerPart.Disposition = disp
					headerPart.DispositionValue = params
				}
				if headerPart.ContentId != "" {
					cid := headerPart.ContentId
					cid = strings.TrimPrefix(cid, "<")
					cid = strings.TrimSuffix(cid, ">")
					headerPart.ContentId = cid

					filename, ok := headerPart.DispositionValue["filename"]
					if ok {
						headerPart.SavedFileName = fmt.Sprintf("%s_%s", cid, filename)
						headerPart.SavedFilePath = fmt.Sprintf("%s/%s_%s", folder, cid, filename)
						err = makeFolder(folder)
						if err != nil {
							emailObj.Error = err.Error()
							emailObjects = append(emailObjects, emailObj)
							continue
						}
						file, err := os.Create(headerPart.SavedFilePath)
						if err != nil {
							emailObj.Error = err.Error()
							emailObjects = append(emailObjects, emailObj)
							continue
						}
						// using io.Copy instead of io.ReadAll to avoid insufficient memory issues
						size, err := io.Copy(file, p.Body)
						if err != nil {
							emailObj.Error = err.Error()
							emailObjects = append(emailObjects, emailObj)
							continue
						}
						headerPart.FileSize = size
						// }

						// for idx, line := range headers {
						// 	//  src="cid:<id>"
						// 	headers[idx].Centent = strings.ReplaceAll(line.Centent, "cid:"+cid, headerPart.SavedFileName)
						// }
					}

				} else {
					b, _ := io.ReadAll(p.Body)
					headerPart.Centent = string(b)
				}

			case *mail.AttachmentHeader:
				contentType, params, err := h.ContentType()
				if err == nil {
					headerPart.ContentType = contentType
					headerPart.ContentTypeValue = params
				}
				disp, params, err := h.ContentDisposition()
				if err == nil {
					headerPart.Disposition = disp
					headerPart.DispositionValue = params
				}

				// This is an attachment
				filename, _ := h.Filename()
				headerPart.Attachment = filename
				headerPart.SavedFileName = filename
				headerPart.SavedFilePath = fmt.Sprintf("%s/%s", folder, filename)
				err = makeFolder(folder)
				if err != nil {
					emailObj.Error = err.Error()
					emailObjects = append(emailObjects, emailObj)
					continue
				}
				file, err := os.Create(headerPart.SavedFilePath)
				if err != nil {
					emailObj.Error = err.Error()
					emailObjects = append(emailObjects, emailObj)
					continue
				}
				// using io.Copy instead of io.ReadAll to avoid insufficient memory issues
				size, err := io.Copy(file, p.Body)
				if err != nil {
					emailObj.Error = err.Error()
					emailObjects = append(emailObjects, emailObj)
					continue
				}
				headerPart.FileSize = size
				emailObj.Attachments = append(emailObj.Attachments, filename)

			}
			headers = append(headers, headerPart)

		}
		emailObj.Body = headers
		emailObjects = append(emailObjects, emailObj)

		// Wait for the update to complete
		// updatedMsg := <-updateCh
		// flags2 := updatedMsg.Uid
		// log.Println("Updated Flags:", flags2)
	}

	if err := <-done; err != nil {
		return nil, err
	}
	// update the item status to seen,next time will not process
	flags := []interface{}{imap.SeenFlag}
	storeItem := imap.FormatFlagsOp(imap.AddFlags, true)
	err = imapClient.Store(seqSet2, storeItem, flags, nil)
	if err != nil {
		return nil, err
	}

	// Close the connection
	if err := imapClient.Close(); err != nil {
		return nil, err
	}

	return emailObjects, nil
}

func receiveImapEmails(email Email) ([]MessageReceived, error) {
	account := email.Account
	
	// 验证账户信息
	if account.Username == "" {
		return nil, fmt.Errorf("IMAP username is required")
	}
	if account.Password == "" {
		return nil, fmt.Errorf("IMAP password is required")
	}
	if account.Server == "" {
		return nil, fmt.Errorf("IMAP server is required")
	}
	if account.Port <= 0 {
		return nil, fmt.Errorf("IMAP port must be greater than 0")
	}

	c, err := connectToServer(account.Username, account.Password, account.Server, account.Port)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	return fetchImapEmails(c, email)

}

// func sendImapEmail(username, password, server string, port int, to, subject, body string) error {
// 	auth := smtp.PlainAuth("", username, password, server)

// 	msg := "From: " + username + "\r\n" +
// 		"To: " + to + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body

// 	err := smtp.SendMail(server+":"+strconv.Itoa(port), auth, username, []string{to}, []byte(msg))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
