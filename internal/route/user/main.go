package user

import (
	"main/db"
	"main/internal/middleware"

	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
    sessionManager middleware.SessionManager
    userRepo db.UserRepo
}

func Init(g *echo.Group, sessionManager middleware.SessionManager, userRepo db.UserRepo) {
    handlers := UserHandlers{
        sessionManager: sessionManager,
        userRepo: userRepo,
    }

    handlers.login(g)
}
