package domain_model

import "fmt"

type DomainError struct {
    Object string
    Code string
}

func (err DomainError) Error() string {
    return fmt.Sprintf("%v error: %v", err.Object, err.Code)
}

func CreateDomainError(object string, code string) DomainError {
    return DomainError {
        Object: object,
        Code: code,
    }
}
