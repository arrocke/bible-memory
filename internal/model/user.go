package model

type User struct {
    Id int
    Email string
    Password string
    FirstName string
    LastName string
}

func (u User) ValidatePassword(password string) bool {
    return password == u.Password
}
