package services

import "fmt"

type NotFoundError struct {
    Resource string
}

func (err NotFoundError) Error() string {
    return fmt.Sprintf("%v not found", err.Resource)
}
