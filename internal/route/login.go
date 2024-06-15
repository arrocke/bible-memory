package route

import (
	"main/internal/middleware"
	"main/view"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type postLoginRequest struct {
    Email string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"required"`
}

func (h Handlers) login(g *echo.Group) {
    g.GET("login", middleware.AuthMiddleware(false)(func(c echo.Context) error {
        page := view.LoginPage(view.LoginPageModel{})
		app := view.Page(view.AppModel{}, page)
		return app.Render(c.Request().Context(), c.Response().Writer)
    }))

    g.POST("login", middleware.AuthMiddleware(false)(func(c echo.Context) error {
        w := c.Response().Writer
        ctx := c.Request().Context()

        var req postLoginRequest
        if err := c.Bind(&req); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        if err := c.Validate(req); err != nil {
            if errors, ok := err.(validator.ValidationErrors); ok {
                model := view.LoginPageModel {
                    Email: req.Email,
                    ValidationErrors: &errors,
                }
                return view.LoginForm(model).Render(ctx, w)
            } else {
                return err
            }
        }

        user, err := h.userRepo.GetByEmail(req.Email)
        if err != nil {
            return err
        }

        if user == nil || !user.ValidatePassword(req.Password) {
            viewModel := view.LoginPageModel{
                Email: req.Email,
                Error: "Invalid email or password.",
            }
            return view.LoginForm(viewModel).Render(ctx, w)
        }

        if err := h.sessionManager.LogIn(c, user.Id()); err != nil {
            return err
        }

        c.Response().Header().Set("Hx-Location", "/passages")
        c.NoContent(http.StatusNoContent)

        return nil
    }))

    g.POST("logout", func(c echo.Context) error {
        if err:= h.sessionManager.LogOut(c); err != nil {
            return err
        }

        c.Response().Header().Set("Hx-Location", "/")
        c.NoContent(http.StatusNoContent)

        return nil
    })
}
