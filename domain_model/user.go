package domain_model

import "math/rand"

type UserProps struct {
    Id int
    Name UserName
    EmailAddress UserEmail
    Password string
}

type User struct {
    props UserProps
    isNew bool
}

func UserFactory(props UserProps) User {
    return User {
        props: props,
        isNew: false,
    }
}


type NewUserProps struct {
    Name UserName
    EmailAddress UserEmail
    Password string
}

func NewUser(props NewUserProps) User {
    id := rand.Int()
    return User {
        props: UserProps {
            Id: id,
            EmailAddress: props.EmailAddress,
            Name: props.Name,
            Password: props.Password,
        },
        isNew: true,
    }
}

func (u *User) IsNew() bool {
    return u.isNew
}

func (u *User) Id() int {
    return u.props.Id
}

func (u *User) Props() UserProps {
    return u.props
}

func (u *User) ValidatePassword(attempt string) bool {
    return u.props.Password == attempt
}

func (u *User) ChangeName(name UserName) {
    u.props.Name = name
}

func (u *User) ChangeEmail(email UserEmail) {
    u.props.EmailAddress = email
}

func (u *User) ChangePassword(newPassword string) {
    u.props.Password = newPassword
}
