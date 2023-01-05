package mail

import (
	"context"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type Mail struct {
	usePlainText      bool
	senderAddress     string
	smtpHostAddr      string
	smtpAuth          smtp.Auth
	receiverAddresses []string
}

func New(senderAddress, smtpHostAddress string) *Mail {
	return &Mail{
		usePlainText:      false,
		senderAddress:     senderAddress,
		smtpHostAddr:      smtpHostAddress,
		receiverAddresses: []string{},
	}
}

type BodyType int

const (
	// PlainText is used to specify that the body is plain text.
	PlainText BodyType = iota
	// HTML is used to specify that the body is HTML.
	HTML
)

func (m *Mail) AuthenticateSMTP(identity, userName, password, host string) {
	m.smtpAuth = smtp.PlainAuth(identity, userName, password, host)
}

func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

func (m *Mail) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		m.usePlainText = true
	default:
		m.usePlainText = false
	}
}

func (m *Mail) newEmail(subject, message string) *email.Email {
	msg := &email.Email{
		To:      m.receiverAddresses,
		From:    m.senderAddress,
		Subject: subject,
		Headers: textproto.MIMEHeader{},
	}

	if m.usePlainText {
		msg.Text = []byte(message)
	} else {
		msg.HTML = []byte(message)
	}
	return msg
}

func (m Mail) Send(ctx context.Context, subject, message string) error {
	msg := m.newEmail(subject, message)

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		err = msg.Send(m.smtpHostAddr, m.smtpAuth)
		if err != nil {
			err = errors.Wrap(err, "failed to send mail")
		}
	}

	return err
}
