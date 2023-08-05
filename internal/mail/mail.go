package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error
}

type SenderCred struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

const (
	gmailSmtpType    = "smtp.gmail.com"
	gmailSmtpAddress = "smtp.gmail.com:587"
	testSmtpType     = "sandbox.smtp.mailtrap.io"
	testSmtpAddress  = "sandbox.smtp.mailtrap.io:25"
)

func NewEmailSender(name, email, password string) EmailSender {
	return &SenderCred{
		name:              name,
		fromEmailAddress:  email,
		fromEmailPassword: password,
	}
}

func (s *SenderCred) SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.name, s.fromEmailAddress)
	e.To = to
	e.Subject = subject
	e.HTML = []byte(content)
	e.Cc = cc
	e.Bcc = bcc

	for _, v := range attachFiles {
		_, err := e.AttachFile(v)

		if err != nil {
			return err
		}
	}

	err := e.Send(testSmtpAddress, smtp.PlainAuth("", s.name, s.fromEmailPassword, testSmtpType))

	return err
}
