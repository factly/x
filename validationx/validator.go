package validationx

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/factly/x/errorx"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Validator - go-playground validator
var Validator validator.Validate

// Trans - Translator
var Trans ut.Translator

// Check - Check struct fields
func Check(model interface{}) []errorx.Message {
	v := validator.New()

	Validator, Trans = addTranslator(v)

	/**
	 **	Register a custom name mapper to be in lower case error message
	 **	Ex- FirstName - first_name
	 **/
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := Validator.Struct(model)

	if err == nil {
		return nil
	}

	errorsList := make([]errorx.Message, 0)
	for _, e := range err.(validator.ValidationErrors) {
		errorsList = append(errorsList, errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Source:  e.Field(),
			Message: e.Translate(Trans),
		})

	}

	return errorsList

}
