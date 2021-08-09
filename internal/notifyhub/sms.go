package notifyhub

import (
	"net/http"
	"time"

	"github.com/taudelta/notifyhub/internal/notifyhub/config"
	"github.com/taudelta/notifyhub/internal/notifyhub/intistelecom"
)

var client intistelecom.Client

func InitSms(cfg config.SmsConfig) {

	smsHttpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client = intistelecom.Client{
		Client: smsHttpClient,
		Login:  cfg.Login,
		Key:    cfg.Key,
	}
}

func SendBySms(notification Notification) error {

	_, err := client.SendSms(notification.Recipient.Phone, notification.Body)
	if err != nil {
		return err
	}

	return nil
}
