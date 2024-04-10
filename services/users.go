package services

import (
	"fmt"
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
    user := domain_model.NewUser(request.FirstName, request.LastName, request.EmailAddress, request.Password)

    if err := service.userRepo.Create(&user); err != nil {
        return 0, nil
    }
    
    return user.Id, nil
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

    user.ChangeName(request.FirstName, request.LastName)
    user.ChangeEmail(request.EmailAddress)
    if (request.Password != "") {
        user.ChangePassword(request.Password)
    }

	return service.userRepo.Update(user)
}

func (service *UserService) ValidatePassword(email string, password string) (*int, error) {
	user, err := service.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

    if user.ValidatePassword(password) {
        return &user.Id, nil
    }

    return nil, nil
}
