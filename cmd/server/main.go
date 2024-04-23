package main

import (
	"context"
	"main/db"
	"main/internal/middleware"
	"main/internal/route/user"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Validator struct {
    validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
  return cv.validator.Struct(i)
}

func main() {
	godotenv.Load(".env")

    e := echo.New()
    e.Validator = &Validator{validator: validator.New()}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
	    e.Logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

    sessionManager := middleware.InitSessions(os.Getenv("SESSION_KEY"))

    e.Use(sessionManager.Middleware())

    g := e.Group("/")

    userRepo := db.CreatePgUserRepo(conn)

    user.Init(g, sessionManager, userRepo)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

    e.Static("/assets", "assets")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
