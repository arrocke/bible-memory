package domain_model

type User struct {
    Id int
    FirstName string
    LastName string
    EmailAddress string
    Password string
}

func NewUser(emailAddress string, firstName string, lastName string, password string) User {
    return User {
        EmailAddress: emailAddress,
        FirstName: firstName,
        LastName: lastName,
        Password: password,
    }
}

func (u *User) ValidatePassword(attempt string) bool {
    return u.Password == attempt
}

func (u *User) ChangeName(firstName string, lastName string) {
    u.FirstName = firstName
    u.LastName = lastName
}

func (u *User) ChangeEmail(email string) {
    u.EmailAddress = email
}

func (u *User) ChangePassword(newPassword string) {
    u.Password = newPassword
}
