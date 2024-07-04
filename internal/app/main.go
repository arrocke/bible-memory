package app

import (
	"context"
	"fmt"
	"main/internal/db"
	"main/internal/view"

	"github.com/benbjohnson/hashfs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfiguration struct {
    Env string
    DatabaseUrl string
    SessionKey string
    Port string
}

type ServerContext struct {
    UserRepo db.UserRepo
    PassageRepo db.PassageRepo
}

func Start(config ServerConfiguration) {
    e := echo.New()
    e.Validator = createValidator()

    view.ShouldHashFiles = config.Env != "development"

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
        PassageRepo: db.PassageRepo{Pool: pool},
    }

    e.Use(SessionMiddleware(config.SessionKey))

    userRoutes(e, context)
    passageRoutes(e, context)

    if view.ShouldHashFiles {
        e.Any("/assets/*", echo.WrapHandler(hashfs.FileServer(view.HashedFS)))
    } else {
        e.Static("/assets",  "internal/view/assets")
    }

    port := config.Port
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
