package domain_model

type UserProps struct {
    Id int
    FirstName string
    LastName string
    EmailAddress string
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
    FirstName string
    LastName string
    EmailAddress string
    Password string
}

func NewUser(props NewUserProps) User {
    return User {
        props: UserProps {
            EmailAddress: props.EmailAddress,
            FirstName: props.FirstName,
            LastName: props.LastName,
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

func (u *User) ChangeName(firstName string, lastName string) {
    u.props.FirstName = firstName
    u.props.LastName = lastName
}

func (u *User) ChangeEmail(email string) {
    u.props.EmailAddress = email
}

func (u *User) ChangePassword(newPassword string) {
    u.props.Password = newPassword
}
