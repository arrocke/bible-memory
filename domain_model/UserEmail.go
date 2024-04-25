package domain_model

import (
	"strings"
)

type UserEmail string

func NewUserEmail(email string) (UserEmail, error) {
    if len(email) == 0 {
        return UserEmail(""), CreateDomainError("UserEmail", "Empty")
    }
    if !strings.Contains(email, "@") {
        return UserEmail(""), CreateDomainError("UserEmail", "InvalidFormat")
    }

    return UserEmail(email), nil
}

func (e UserEmail) Value() string {
	return string(e)
}
