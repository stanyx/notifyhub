package recipients

import (
	"encoding/json"
	"net/http"

	log "github.com/taudelta/nanolog"
	"github.com/taudelta/notifyhub/app/web"
)

type Recipient struct {
	Phone    string
	Email    string
	Telegram string
}

type RecipientResponse struct {
	Id       int    `json:"id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
}

func GetRecipients(w http.ResponseWriter, r *http.Request) {

	response, err := GetRecipientsFromDB()
	if err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "get_recipients_error", "")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
	if err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "get_recipients_error", "")
	}

}

type RecipientCreateResponse struct {
	Id int64 `json:"id"`
}

func CreateRecipient(w http.ResponseWriter, r *http.Request) {

	var form Recipient

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "create_recipient_error", "")
		return
	}

	id, err := CreateRecipientDB(form)
	if err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "create_recipient_error", "")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&RecipientCreateResponse{
		Id: id,
	})

	if err != nil {
		log.Error().Println(err)
		web.JsonError(w, 500, "create_recipient_error", "")
	}

}
