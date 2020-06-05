package validationx

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Validator - go-playground validator
var Validator validator.Validate

// Trans - Translator
var Trans ut.Translator

// Check - Check struct fields
func Check(model interface{}) interface{} {
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

	var arr []interface{}
	for _, e := range err.(validator.ValidationErrors) {
		arr = append(arr, map[string]string{
			"field":   e.Field(),
			"message": e.Translate(Trans),
		})

	}

	return arr
}
