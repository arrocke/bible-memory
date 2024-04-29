package user

import (
	"main/domain_model"
	"main/view"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type postRegisterRequest struct {
    FirstName string `form:"first_name" validate:"required"`
    LastName string `form:"last_name" validate:"required"`
    Email string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"required,min=8"`
    ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
}

func (h Handlers) register(g *echo.Group) {
    g.GET("register", func (c echo.Context) error {
        page := view.RegisterPage(view.RegisterPageModel{})
        app := view.Page(view.AppModel{}, page)
        return app.Render(c.Request().Context(), c.Response().Writer)
    })

    g.POST("register", func (c echo.Context) error {
        w := c.Response().Writer
        ctx := c.Request().Context()

        var req postRegisterRequest
        if err := c.Bind(&req); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        if err := c.Validate(req); err != nil {
            if errors, ok := err.(validator.ValidationErrors); ok {
                model := view.RegisterPageModel {
                    FirstName: req.FirstName,
                    LastName: req.LastName,
                    Email: req.Email,
                    ValidationErrors: &errors,
                }
                return view.RegisterForm(model).Render(ctx, w)
            } else {
                return err
            }
        }

        name, err := domain_model.NewUserName(req.FirstName, req.LastName)
        if err != nil {
            return err
        }
        email, err := domain_model.NewUserEmail(req.Email)
        if err != nil {
            return err
        }
        user := domain_model.NewUser(domain_model.NewUserProps{
            Name: name,
            EmailAddress: email,
            Password: req.Password,
        })

        if err := h.userRepo.Commit(user); err != nil {
            return err
        }

        
        if err := h.sessionManager.LogIn(c, user.Id()); err != nil {
            return err
        }

        c.Response().Header().Set("Hx-Location", "/passages")
        c.NoContent(http.StatusNoContent)

        return nil
    })
}
