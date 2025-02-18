package mailer

import (
	"fmt"
	"github.com/go-mail/mail/v2"
	"time"
)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, userName string, token string) error {

	content := fmt.Sprintf("welcome %s you are a few steps aways from activating your account"+
		"your activation code is %s", userName, token)
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", "Welcome Mail")
	msg.SetBody("text/plain", content)
	err := m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
