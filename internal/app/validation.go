package app

import (
	"main/internal/model"

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

func validateReference(fl validator.FieldLevel) bool {
    return model.ReferenceFormat.Match([]byte(fl.Field().String()))
}

