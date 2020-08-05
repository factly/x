package errorx

import (
	"net/http"

	"github.com/factly/x/renderx"
)

// Message - error model
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

// InvalidID error
func InvalidID() Message {
	return Message{
		Code:    http.StatusNotFound,
		Message: "Invalid ID",
	}
}

// InternalServerError error
func InternalServerError() Message {
	return Message{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
	}
}

// DBError - Database errors
func DBError() Message {
	return Message{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong with db queries",
	}
}

//RecordNotFound - record not found error
func RecordNotFound() Message {
	return Message{
		Code:    http.StatusNotFound,
		Message: "Record not found",
	}
}

//DecodeError - errors while decoding request body
func DecodeError() Message {
	return Message{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid request body",
	}
}

//NetworkError - errors while connection refused
func NetworkError() Message {
	return Message{
		Code:    http.StatusServiceUnavailable,
		Message: "Connection failed",
	}
}

//CannotSaveChanges - errors when an item cannot be changed
func CannotSaveChanges() Message {
	return Message{
		Code:    http.StatusUnprocessableEntity,
		Message: "Can not save changes",
	}
}

type response struct {
	Errors []Message `json:"errors"`
}

// Render -  render errors
func Render(w http.ResponseWriter, err []Message) {

	result := response{
		Errors: err,
	}
	renderx.JSON(w, err[0].Code, result)
}

// Parser - to parse error Messages
func Parser(msg Message) []Message {

	err := make([]Message, 1)
	err[0] = msg
	return err
}
