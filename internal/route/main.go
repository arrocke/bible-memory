package route

import (
	"main/db"
	"main/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
    sessionManager middleware.SessionManager
    dbConn *pgxpool.Pool
    userRepo db.UserRepo
    passageRepo db.PassageRepo
}

func Init(g *echo.Group, sessionManager middleware.SessionManager, dbConn *pgxpool.Pool, userRepo db.UserRepo, passageRepo db.PassageRepo) {
    handlers := Handlers{
        sessionManager: sessionManager,
        dbConn: dbConn,
        userRepo: userRepo,
        passageRepo: passageRepo,
    }

    handlers.login(g)
    handlers.register(g)
    handlers.profile(g)
    handlers.passages(g)
}
