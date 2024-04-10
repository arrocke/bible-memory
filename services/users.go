package services

import (
	"main/db"
	"main/domain_model"
)

type UserService struct {
	userRepo db.UserRepo
}

func CreateUsersService(userRepo db.UserRepo) UserService {
    return UserService{userRepo}
}

type CreateUserRequest struct {
    FirstName string
    LastName string
    EmailAddress string
    Password string
}

func (service *UserService) Create(request CreateUserRequest) (int, error) {
    name, err := domain_model.NewUserName(request.FirstName, request.LastName)
    if err != nil {
        return 0, err
    }

    user := domain_model.NewUser(domain_model.NewUserProps{
        Name: name,
        EmailAddress: request.EmailAddress,
        Password: request.Password,
    })

    if err := service.userRepo.Commit(&user); err != nil {
        return 0, nil
    }
    
    return user.Id(), nil
}

type UpdateProfileRequest struct {
    Id int
    FirstName string
    LastName string
    EmailAddress string
    Password string
}

func (service *UserService) UpdateProfile(request UpdateProfileRequest) error {
	user, err := service.userRepo.Get(request.Id)
	if err != nil {
		return err
	}

    name, err := domain_model.NewUserName(request.FirstName, request.LastName)
    if err != nil {
        return err
    }

    user.ChangeName(name)
    user.ChangeEmail(request.EmailAddress)
    if (request.Password != "") {
        user.ChangePassword(request.Password)
    }

	return service.userRepo.Commit(user)
}

func (service *UserService) ValidatePassword(email string, password string) (*int, error) {
	user, err := service.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

    id := user.Id()
    if user.ValidatePassword(password) {
        return &id, nil
    }

    return nil, nil
}
