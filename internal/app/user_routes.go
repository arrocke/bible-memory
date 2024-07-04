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
        return RenderComponent(c, view.Html(loginView))
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

    e.POST("/logout", func(c echo.Context) error {
        if err := LogOut(c); err != nil {
            return err
        }

        return RedirectWithRefresh(c, "/")
    })

    e.GET("/profile", func(c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        user, err := ctx.UserRepo.GetUserById(c.Request().Context(), userId)
        if err != nil {
            return err
        }

        component := view.ProfileView(view.ProfileViewModel{
            Email: user.Email,
            FirstName: user.FirstName,
            LastName: user.LastName,
        })
        if IsHtmxRequest(c) {
            Retarget(c, "#view")
            return RenderComponent(c, component)
        } else {
            page, err := ctx.CreateView(c, component)
            if err != nil {
                return err
            }
            return RenderComponent(c, view.Html(page))
        }
    }, AuthMiddleware(true))

    type putProfileRequest struct {
        FirstName string `form:"first_name" validate:"required"`
        LastName string `form:"last_name" validate:"required"`
        Email string `form:"email" validate:"required,email"`
    }

    e.PUT("/profile", func(c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

		var req putProfileRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				model := view.ProfileViewModel{
                    FirstName: req.FirstName,
                    LastName: req.LastName,
                    Email: req.Email,
					Errors: &errors,
				}
				return RenderComponent(c, view.ProfileForm(model))
			} else {
				return err
			}
		}

        user, err := ctx.UserRepo.GetUserById(c.Request().Context(), userId)
        if err != nil {
            return err
        }

        user.Email = req.Email
        user.FirstName = req.FirstName
        user.LastName = req.LastName

        if err := ctx.UserRepo.Update(c.Request().Context(), user); err != nil {
            return err
        }

        model := view.ProfileViewModel{
            FirstName: req.FirstName,
            LastName: req.LastName,
            Email: req.Email,
            Success: true,
        }
        return RenderComponent(c, view.ProfileForm(model))
    }, AuthMiddleware(true))
}
