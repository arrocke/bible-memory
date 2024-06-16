package app

import (
	"context"
	"fmt"
	"main/internal/db"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfiguration struct {
    DatabaseUrl string
    SessionKey string
    Port string
}

type Validator struct {
    validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
  return cv.validator.Struct(i)
}

type ServerContext struct {
    UserRepo db.UserRepo
}

func Start(config ServerConfiguration) {
    e := echo.New()
    e.Validator = &Validator{validator: validator.New()}

    pool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
	    e.Logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

    e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
        LogStatus: true,
        LogURI:    true,
        LogMethod: true,
        LogError: true,
        LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
            if v.Error != nil {
                fmt.Printf("REQUEST: %v %v %v %v\n", v.Method, v.URI, v.Status, v.Error.Error())
            } else {
                fmt.Printf("REQUEST: %v %v %v\n", v.Method, v.URI, v.Status)
            }
            return nil
        },
    }))

    context := ServerContext{
        UserRepo: db.UserRepo{Pool: pool},
    }

    e.Use(SessionMiddleware(config.SessionKey))

    userRoutes(e, context)

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Passages")
    })

    e.Static("/assets", "assets")


    port := config.Port
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
