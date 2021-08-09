package notifications

import (
	"time"

	"github.com/taudelta/notifyhub/internal/notifyhub/db"
)

func CreateRecipientInfo(recipient db.Recipient) (int64, error) {

	var recipientInfoID int64

	connection := db.Connection()

	err := connection.QueryRow(
		"insert into notification_recipient_info (email, phone, telegram) "+
			"values ($1, $2, $3) returning id",
		recipient.Email, recipient.Phone, recipient.Telegram,
	).Scan(&recipientInfoID)
	if err != nil {
		return 0, nil
	}

	return recipientInfoID, nil
}

func CreateNotificationDB(recipientInfoID int64, notificationType, message string) (int64, error) {

	var notificationID int64

	connection := db.Connection()

	err := connection.QueryRow(
		"insert into notification (status, timestamp, recipient_id, message, notification_type) "+
			"values ('processing', $1, $2, $3, $4) returning id",
		time.Now(), recipientInfoID, message, notificationType,
	).Scan(&notificationID)
	if err != nil {
		return 0, err
	}

	return notificationID, nil
}
