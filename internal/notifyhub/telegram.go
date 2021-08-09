package notifyhub

import (
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/taudelta/nanolog"
)

var bot *tgbotapi.BotAPI

func InitTelegram(token string) error {

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	var err error
	bot, err = tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		return err
	}

	return nil
}

func SendByTelegram(notification Notification) error {

	chatId, err := strconv.ParseInt(notification.Recipient.TelegramChatId, 10, 64)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, notification.Body)
	if _, err := bot.Send(msg); err != nil {
		log.Error().Println("send message error: %s", err)
		return err
	}

	return nil
}
