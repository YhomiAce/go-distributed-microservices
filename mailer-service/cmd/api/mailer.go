package main

import (
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain  string
	Host   string
	Port	int
	Username	string
	Password	string
	Encryption	string
	FromAddress	string
	FromName	string
}

type Message struct {
	From	string
	FromAddress	string
	To		string
	Subject	string
	Attachments []string
	Data 	any
	DataMap map[string]any
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
		case "tls":
			return mail.EncryptionSTARTTLS
		case "ssl":
			return mail.EncryptionSSLTLS
		case "none":
			return  mail.EncryptionNone
		default:
			return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) sendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromName
	}
	if msg.FromAddress == "" {
		msg.FromAddress = m.FromAddress
	}

	plainMessage, ok := msg.Data.(string)
	if !ok {
		plainMessage = "Invalid message content"
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 30 * time.Second
	server.SendTimeout = 30 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println("failed to connect",err.Error())
		return  err;
	}

	email :=mail.NewMSG()
	email.SetFrom(msg.FromAddress).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)

	htmlMessage := "<p>" + plainMessage + "</p>"

	email.AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _,file := range msg.Attachments {
			email.AddAttachment(file)
		}
	}
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}