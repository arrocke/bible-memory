package domain_model

type UserName struct {
    firstName string
    lastName string
}

func NewUserName(firstName string, lastName string) (UserName, error) {
    if len(firstName) == 0 {
        return UserName{}, CreateDomainError("UserName", "FirstNameEmpty")
    }
    if len(lastName) == 0 {
        return UserName{}, CreateDomainError("UserName", "LastNameEmpty")
    }

    return UserName{firstName: firstName, lastName: lastName}, nil
}

func (n UserName) FirstName() string {
    return n.firstName
}
func (n UserName) LastName() string {
    return n.lastName
}
