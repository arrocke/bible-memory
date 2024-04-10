package domain_model

import "fmt"

type UserName struct {
    firstName string
    lastName string
}

func NewUserName(firstName string, lastName string) (UserName, error) {
    if len(firstName) == 0 || len(lastName) == 0 {
        return UserName{},fmt.Errorf("UserName requires first and last name")
    }

    return UserName{firstName: firstName, lastName: lastName}, nil
}

func (n UserName) FirstName() string {
    return n.firstName
}
func (n UserName) LastName() string {
    return n.lastName
}
