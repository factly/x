package errorx

import (
	"net/http"
	"testing"
)

func TestErrorx(t *testing.T) {

	t.Run("Parser with a message", func(t *testing.T) {
		mess := Message{
			Code:    404,
			Message: "Not Found",
			Source:  "file/c.go",
		}

		retVal := Parser(mess)

		if len(retVal) != 1 {
			t.Errorf("Returned object must contain one element")
		}

		if retVal[0].Code != mess.Code {
			t.Errorf("Returned wrong code expected %v, got %v", mess.Code, retVal[0].Code)
		}

		if retVal[0].Message != mess.Message {
			t.Errorf("Returned wrong message expected %v, got %v", mess.Message, retVal[0].Message)
		}

		if retVal[0].Source != mess.Source {
			t.Errorf("Returned wrong source expected %v, got %v", mess.Source, retVal[0].Source)
		}
	})

	t.Run("get InvalidID error message", func(t *testing.T) {
		err := InvalidID()

		if err.Code != http.StatusBadRequest {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusNotFound, err.Code)
		}

		if err.Message != "Invalid ID" {
			t.Errorf("Returned wrong message expected %v, got %v", "Invalid ID", err.Message)
		}
	})

	t.Run("get InternalServerError error message", func(t *testing.T) {
		err := InternalServerError()

		if err.Code != http.StatusInternalServerError {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusInternalServerError, err.Code)
		}

		if err.Message != "Something went wrong" {
			t.Errorf("Returned wrong message expected %v, got %v", "Something went wrong", err.Message)
		}
	})

	t.Run("get DBError error message", func(t *testing.T) {
		err := DBError()

		if err.Code != http.StatusInternalServerError {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusInternalServerError, err.Code)
		}

		if err.Message != "Something went wrong with db queries" {
			t.Errorf("Returned wrong message expected %v, got %v", "Something went wrong with db queries", err.Message)
		}
	})

	t.Run("get RecordNotFound error message", func(t *testing.T) {
		err := RecordNotFound()

		if err.Code != http.StatusNotFound {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusNotFound, err.Code)
		}

		if err.Message != "Record not found" {
			t.Errorf("Returned wrong message expected %v, got %v", "Record not found", err.Message)
		}
	})

	t.Run("get DecodeError error message", func(t *testing.T) {
		err := DecodeError()

		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusUnprocessableEntity, err.Code)
		}

		if err.Message != "Invalid request body" {
			t.Errorf("Returned wrong message expected %v, got %v", "Invalid request body", err.Message)
		}
	})

	t.Run("get NetworkError error message", func(t *testing.T) {
		err := NetworkError()

		if err.Code != http.StatusServiceUnavailable {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusServiceUnavailable, err.Code)
		}

		if err.Message != "Connection failed" {
			t.Errorf("Returned wrong message expected %v, got %v", "Connection failed", err.Message)
		}
	})

	t.Run("get CannotSaveChanges error message", func(t *testing.T) {
		err := CannotSaveChanges()

		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusUnprocessableEntity, err.Code)
		}

		if err.Message != "Can not save changes" {
			t.Errorf("Returned wrong message expected %v, got %v", "Can not save changes", err.Message)
		}
	})

	t.Run("get Unauthorized error message", func(t *testing.T) {
		err := Unauthorized()

		if err.Code != http.StatusUnauthorized {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusUnauthorized, err.Code)
		}

		if err.Message != "Not allowed" {
			t.Errorf("Returned wrong message expected %v, got %v", "Not allowed", err.Message)
		}
	})

	t.Run("get SameNameExist error message", func(t *testing.T) {
		err := SameNameExist()

		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusUnprocessableEntity, err.Code)
		}

		if err.Message != "Entity with same name exists" {
			t.Errorf("Returned wrong message expected %v, got %v", "Entity with same name exists", err.Message)
		}
	})

	t.Run("get CannotDelete error message", func(t *testing.T) {
		err := CannotDelete("tag", "post")

		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong code expected %v, got %v", http.StatusUnprocessableEntity, err.Code)
		}

		if err.Message != "tag is associated with some post" {
			t.Errorf("Returned wrong message expected %v, got %v", "tag is associated with some post", err.Message)
		}
	})
}
