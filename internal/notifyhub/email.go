package notifyhub

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/taudelta/notifyhub/internal/notifyhub/config"
)

func NewEmailMessage(recipients []string, subject, body string, attachments [][]byte) (*email.Email, error) {

	mail := &email.Email{
		To:      recipients,
		From:    "thinkandwin5000@gmail.com",
		Subject: subject,
		Text:    []byte(body),
	}

	for _, attachment := range attachments {
		_, err := mail.Attach(bytes.NewReader(attachment), "attachment", "")
		if err != nil {
			return nil, err
		}
	}

	return mail, nil
}

var emailConfig config.EmailConfig

func SetEmailConfig(cfg config.EmailConfig) {
	emailConfig = cfg
}

func SendByEmail(notification Notification) error {

	addr := fmt.Sprintf("%s:%d", emailConfig.Host, emailConfig.Port)
	auth := smtp.PlainAuth("", emailConfig.Login, emailConfig.Password, emailConfig.Host)

	message, err := NewEmailMessage(
		[]string{"backupchik1@gmail.com"},
		"NotifyHub. Notification",
		"message",
		[][]byte{},
	)
	if err != nil {
		return err
	}

	if emailConfig.UseTLS {

		tlsConfig := &tls.Config{
			ServerName: emailConfig.Host,
		}

		if emailConfig.InsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
		} else {
			cer, err := tls.LoadX509KeyPair(emailConfig.TLSCrtFile, emailConfig.TLSKeyFile)
			if err != nil {
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{cer}
		}
		return message.SendWithTLS(addr, auth, tlsConfig)
	}

	return message.Send(addr, auth)
}
