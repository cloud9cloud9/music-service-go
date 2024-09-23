package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var Validate = validator.New()

func ParseJSON(request *http.Request, input any) error {
	if request.Body == nil {
		return fmt.Errorf("Missing input body")
	}
	return json.NewDecoder(request.Body).Decode(input)
}

func WriteJSON(writer http.ResponseWriter, status int, input any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(input)
}

func WriteError(writer http.ResponseWriter, status int, err error) error {
	return WriteJSON(writer, status, map[string]string{"error": err.Error()})
}

func InvalidParsingJSON(writer http.ResponseWriter) error {
	return WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid parsing JSON"))
}
