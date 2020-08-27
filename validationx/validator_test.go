package validationx

import (
	"net/http"
	"testing"
)

func TestValidationx(t *testing.T) {

	type TestObj struct {
		Field1 string `json:"field-1" validate:"required"`
		Email  string `validate:"email"`
	}

	t.Run("Validate a object", func(t *testing.T) {
		obj := TestObj{
			Field1: "test field",
			Email:  "test@mail.com",
		}

		errs := Check(obj)

		if len(errs) != 0 {
			t.Errorf("Returned unexpected errors")
		}
	})

	t.Run("passing invalid email", func(t *testing.T) {
		obj := TestObj{
			Field1: "test field",
			Email:  "test@maiom",
		}

		errs := Check(obj)

		if len(errs) != 1 {
			t.Errorf("Returned no errors")
		}

		if errs[0].Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong error code expected %v, got %v", http.StatusUnprocessableEntity, errs[0].Code)
		}

		if errs[0].Message != "Email must be a valid email" {
			t.Errorf("Returned wrong error message expected %v, got %v", "Email must be a valid email", errs[0].Message)
		}
	})

	t.Run("passing empty string in required field", func(t *testing.T) {
		obj := TestObj{
			Field1: "",
			Email:  "test@mail.com",
		}

		errs := Check(obj)

		if len(errs) != 1 {
			t.Errorf("Returned no errors")
		}

		if errs[0].Code != http.StatusUnprocessableEntity {
			t.Errorf("Returned wrong error code expected %v, got %v", http.StatusUnprocessableEntity, errs[0].Code)
		}

		if errs[0].Message != "Field1 is a required field" {
			t.Errorf("Returned wrong error message expected %v, got %v", "Field1 is a required field", errs[0].Message)
		}
	})

	t.Run("passing empty object", func(t *testing.T) {
		obj := TestObj{}

		errs := Check(obj)

		if len(errs) != 2 {
			t.Errorf("Not returned expected errors")
		}

		errMessages := []string{"Field1 is a required field", "Email must be a valid email"}

		for i, e := range errs {
			if e.Code != http.StatusUnprocessableEntity {
				t.Errorf("Returned wrong error code expected %v, got %v", http.StatusUnprocessableEntity, e.Code)
			}

			if e.Message != errMessages[i] {
				t.Errorf("Returned wrong error message expected %v, got %v", "Field1 is a required field", errs[0].Message)
			}
		}
	})
}
