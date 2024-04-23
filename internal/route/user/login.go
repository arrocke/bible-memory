package user

import (
	"main/view"
	"net/http"

	"github.com/labstack/echo/v4"
    "github.com/go-playground/validator/v10"
)

type postLoginRequest struct {
    Email string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"required"`
}

func (h UserHandlers) login(g *echo.Group) {
    g.GET("login", func(c echo.Context) error {
        page := view.LoginPage(view.LoginPageModel{})
		app := view.Page(view.AppModel{}, page)
		return app.Render(c.Request().Context(), c.Response().Writer)
    })

    g.POST("login", func(c echo.Context) error {
        w := c.Response().Writer
        ctx := c.Request().Context()

        var r postLoginRequest
        if err := c.Bind(&r); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        if err := c.Validate(r); err != nil {
            if errors, ok := err.(validator.ValidationErrors); ok {
                model := view.LoginPageModel {
                    Email: r.Email,
                    ValidationErrors: &errors,
                }
                return view.LoginForm(model).Render(ctx, w)
            } else {
                return err
            }
        }

        user, err := h.userRepo.GetByEmail(r.Email)
        if err != nil {
            return err
        }

        if user == nil || !user.ValidatePassword(r.Password) {
            viewModel := view.LoginPageModel{
                Email: r.Email,
                Error: "Invalid email or password.",
            }
            return view.LoginForm(viewModel).Render(ctx, w)
        }

		w.Header().Set("Hx-Location", "/passages")
        w.WriteHeader(http.StatusNoContent)

        return nil
    })
}
