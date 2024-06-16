package model

import "fmt"

type User struct {
    Id int
    Email string
    Password string
    FirstName string
    LastName string
}

func (u User) ValidatePassword(password string) bool {
    fmt.Printf("%v %v", password, u.Password)
    return password == u.Password
}
