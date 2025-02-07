package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var (
	errorNoBody = errors.New("NO_BODY_PROVIDED")
	errorInvalidBody = errors.New("INVALID_BODY")
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
		return errorInvalidBody
	}
	if err := validate.Struct(payload); err != nil{
		return errorInvalidBody
	}
	return nil
}





func ParseAndReadMultiPartForm(w http.ResponseWriter,r *http.Request, maxSize int64, payload any, fieldName string) error{
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	err := r.ParseMultipartForm(maxSize)
	if err != nil{
		return ErrFileTooBig
	}

	text := r.FormValue(fieldName)
	if text == ""{
		return errorNoBody
	}
	if err := json.Unmarshal([]byte(text), &payload); err != nil{
		return errorInvalidBody
	}
	if err := validate.Struct(payload); err != nil{
		return errorInvalidBody
	}
	return nil
} 



func GetErrorCode(err error, errors map[error]int, fallback int) int{
	code, ok := errors[err]
	if !ok{
		return 500
	}
	return code
}