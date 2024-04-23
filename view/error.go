package view

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func defaultValidationMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required": return fmt.Sprintf("%v is required", err.Field())
    case "email": return "Must be a valid email"
    default: return fmt.Sprintf("unhandled error: %v", err.Tag())
    }
}

type validationErrorModel struct {
    err *validator.ValidationErrors
    field string
    messages map[string]string
}

func formatValidationError(model validationErrorModel) string {
    if model.err == nil {
        return ""
    } else {
        for _, err := range *model.err {
            if err.Field() == model.field {
                if msg, ok := model.messages[err.Tag()]; ok {
                    return strings.ReplaceAll(msg, "{param}", err.Param())
                } else {
                    return defaultValidationMessage(err)
                }
            }
        }
        return ""
    }
}

func hasValidationError(model validationErrorModel) bool {
    if model.err == nil {
        return false
    } else {
        for _, err := range *model.err {
            if err.Field() == model.field {
                return true
            }
        }
        return false
    }
}

