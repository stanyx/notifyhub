package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/taudelta/nanolog"
	"github.com/taudelta/notifyhub/app/recipients"
	"github.com/taudelta/notifyhub/app/web"
	"github.com/taudelta/notifyhub/internal/notifyhub"
	"github.com/taudelta/notifyhub/internal/notifyhub/db"
	"github.com/taudelta/umq"
)

type NotificationSend struct {
	RecipientList     []int64  `json:"recipientList"`
	NotificationTypes []string `json:"notificationTypes"`
	Message           string   `json:"message"`
}

type SendError struct {
	RecipientID      int64  `json:"recipient_id"`
	NotificationType string `json:"notificationType"`
	Error            string `json:"error"`
}

type SendResult struct {
	RecipientID      int64  `json:"recipient_id"`
	NotificationType string `json:"notification_type"`
}

type SendResponse struct {
	Enqueued []SendResult `json:"enqueued"`
	Errors   []SendError  `json:"errors"`
}

func SendNotification(messageProducer *umq.RabbitProducer, w http.ResponseWriter, r *http.Request) {

	var form NotificationSend

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "send_error", "")
		return
	}

	var success []SendResult
	var errors []SendError

	for _, recipientID := range form.RecipientList {

		recipient, err := recipients.GetRecipientByID(recipientID)
		if err != nil {
			log.Error().Println(err)
			errorText := "internal error"
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				errorText = "not found"
			}
			errors = append(errors, SendError{
				RecipientID: recipientID,
				Error:       errorText,
			})
			continue
		}

		recipientInfoID, err := CreateRecipientInfo(
			db.Recipient{
				Email:    recipient.Email,
				Phone:    recipient.Phone,
				Telegram: recipient.Telegram,
			},
		)
		if err != nil {
			log.Error().Println(err)
			errors = append(errors, SendError{
				RecipientID: recipientID,
				Error:       "internal_error",
			})
			continue
		}

		for _, notificationType := range form.NotificationTypes {

			notificationID, err := CreateNotificationDB(recipientInfoID, notificationType, form.Message)
			if err != nil {
				log.Error().Println(err)
				errors = append(errors, SendError{
					RecipientID: recipientID,
					Error:       "internal_error",
				})
				continue
			}

			err = notifyhub.SendToQueue(messageProducer, notifyhub.Notification{
				Type: notificationType,
				Recipient: notifyhub.RecipientInfo{
					Email:          recipient.Email,
					Phone:          recipient.Phone,
					TelegramChatId: recipient.Telegram,
				},
				Body: form.Message,
				DbID: notificationID,
			})

			if err != nil {
				errors = append(errors, SendError{
					RecipientID:      recipientID,
					NotificationType: notificationType,
				})
			} else {
				success = append(success, SendResult{
					RecipientID:      recipientID,
					NotificationType: notificationType,
				})
			}
		}
	}

	response := SendResponse{
		Enqueued: success,
		Errors:   errors,
	}

	responseBody, err := json.Marshal(&response)
	if err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "send_error", "")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(responseBody))
}
