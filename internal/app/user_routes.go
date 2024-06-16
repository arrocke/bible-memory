package app

import (
	"errors"
	"main/internal/db"
	"main/internal/view"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type postLoginRequest struct {
    Email string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"required"`
}

func userRoutes(e *echo.Echo, ctx ServerContext) {
    e.GET("/login", func(c echo.Context) error {
        loginView := view.LoginView(view.LoginViewModel{})
        return view.Html(loginView).Render(c.Request().Context(), c.Response().Writer)
    }, AuthMiddleware(false))

    e.POST("/login", func(c echo.Context) error {
        var req postLoginRequest
        if err := c.Bind(&req); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        if err := c.Validate(req); err != nil {
            if errors, ok := err.(validator.ValidationErrors); ok {
                model := view.LoginViewModel {
                    Email: req.Email,
                    Errors: &errors,
                }
                return view.LoginForm(model).Render(c.Request().Context(), c.Response().Writer)
            } else {
                return err
            }
        }

        user, err := ctx.UserRepo.GetUserByEmail(c.Request().Context(), req.Email)
        userNotFound := false
        if err != nil {
            if errors.Is(err, db.NotFoundError) {
                userNotFound = true
            } else {
                return err
            }
        }

        if userNotFound || !user.ValidatePassword(req.Password) {
            model := view.LoginViewModel {
                Email: req.Email,
                Error: "Invalid email or password.",
            }
            return view.LoginForm(model).Render(c.Request().Context(), c.Response().Writer)
        }

        if err := LogIn(c, user.Id); err != nil {
            return err
        }
        
        return RedirectWithRefresh(c, "/")
    }, AuthMiddleware(false))
}
