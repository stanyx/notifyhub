package notifyhub

import (
	"encoding/json"

	log "github.com/taudelta/nanolog"
	"github.com/taudelta/umq"
)

type RecipientInfo struct {
	Email          string
	Phone          string
	TelegramChatId string
}

type Notification struct {
	Type       string
	TemplateID string
	Recipient  RecipientInfo
	Body       string
	Context    map[string]interface{}
	DbID       int64
}

func ParseMessage(body []byte) (*Notification, error) {
	var notification Notification
	err := json.Unmarshal(body, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func SendToQueue(publisher *umq.RabbitProducer, notification Notification) error {

	// send message to rabbitmq

	body, err := json.Marshal(&notification)
	if err != nil {
		return err
	}

	return publisher.Send(body)

}

func SendMessage(notification Notification) error {

	log.Debug().Println("send notification", notification)

	switch notification.Type {
	case "email":
		return SendByEmail(notification)
	case "sms":
		return SendBySms(notification)
	case "telegram":
		return SendByTelegram(notification)
	default:
		log.Error().Printf("unsupported notification type: %s", notification.Type)
	}

	return nil
}
