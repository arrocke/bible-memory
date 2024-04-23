package user

import (
	"main/db"

	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
    userRepo db.UserRepo
}

func Init(g *echo.Group, userRepo db.UserRepo) {
    handlers := UserHandlers{
        userRepo: userRepo,
    }

    handlers.login(g)
}
