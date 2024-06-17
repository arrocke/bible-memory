package app

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type requestValidator struct {
    validator *validator.Validate
}

func (cv *requestValidator) Validate(i interface{}) error {
  return cv.validator.Struct(i)
}

func createValidator() *requestValidator {
    v := requestValidator{validator: validator.New()}

    v.validator.RegisterValidation("reference", validateReference)

    return &v
}

var referenceRegexp = regexp.MustCompile(`(.+?)\s*(\d+)[.:](\d+)(?:\s*-\s*(?:(\d+)[.:])?(\d+))?`)

func validateReference(fl validator.FieldLevel) bool {
    return referenceRegexp.Match([]byte(fl.Field().String()))
}

