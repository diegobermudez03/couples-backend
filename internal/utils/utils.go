package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var (
	errorNoBody = errors.New("no body provided")
	errorInvalidBody = errors.New("invalid body")
)

func WriteJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, status int, err error) error{
	return WriteJSON(
		w,
		status, 
		map[string]any{
			"error" : err.Error(),
		},
	)
}



func ReadJSON(r *http.Request, payload any) error{
	if r.Body == nil{
		return errorNoBody
	}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil{
		log.Println(err.Error())
		return errorInvalidBody
	}
	if err := validate.Struct(payload); err != nil{
		return errorInvalidBody
	}
	return nil
}