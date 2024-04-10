package domain_model

import (
	"fmt"
	"strings"
)

type UserEmail string

func NewUserEmail(email string) (UserEmail, error) {
    if len(email) == 0 {
        return UserEmail(""),fmt.Errorf("email is required")
    }
    if !strings.Contains(email, "@") {
        return UserEmail(""),fmt.Errorf("invalid email format")
    }

    return UserEmail(email),nil
}

func (e UserEmail) Value() string {
	return string(e)
}
