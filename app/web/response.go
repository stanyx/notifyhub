package web

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func Error(w http.ResponseWriter, statusCode int, contentType string, errCode, errDescription string) {

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)

	encoder := json.NewEncoder(w)
	encoder.Encode(&ErrorResponse{
		Code:        errCode,
		Description: errDescription,
	})
}

func JsonError(w http.ResponseWriter, statusCode int, errCode, errDescription string) {
	Error(w, statusCode, "application/json", errCode, errDescription)
}
