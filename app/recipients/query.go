package recipients

import (
	"github.com/taudelta/notifyhub/internal/notifyhub/db"
)

func GetRecipientByID(recipientId int64) (*RecipientResponse, error) {
	conn := db.Connection()

	var row RecipientResponse

	err := conn.QueryRow(
		"select id, phone, email, telegram from recipient where id = $1",
		recipientId,
	).Scan(&row.Id, &row.Phone, &row.Email, &row.Telegram)

	if err != nil {
		return nil, err
	}

	return &row, nil
}

func GetRecipientsFromDB() ([]RecipientResponse, error) {

	conn := db.Connection()
	rows, err := conn.Query(
		"select id, phone, email, telegram from recipient",
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var response []RecipientResponse

	for rows.Next() {
		var row RecipientResponse
		if err := rows.Scan(&row.Id, &row.Phone, &row.Email, &row.Telegram); err != nil {
			return nil, err
		}
		response = append(response, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

func CreateRecipientDB(form Recipient) (int64, error) {
	var id int64
	conn := db.Connection()
	err := conn.QueryRow(
		"insert into recipient (phone, email, telegram) "+
			"values ($1, $2, $3) returning id",
		form.Phone, form.Email, form.Telegram,
	).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}
