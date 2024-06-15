package route

import (
	"main/domain_model"
	"main/internal/middleware"
	"main/view"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type postProfileRequest struct {
    FirstName string `form:"first_name" validate:"required"`
    LastName string `form:"last_name" validate:"required"`
    Email string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"omitempty,min=8"`
    ConfirmPassword string `form:"confirm_password" validate:"eqfield=Password"`
}

func (h Handlers) profile(g *echo.Group) {
    g.GET("profile", middleware.AuthMiddleware(true)(func(c echo.Context) error {
        userId := c.Get("user_id").(int)

        user, err := h.userRepo.Get(userId)
        if err != nil {
            return err
        }

        props := user.Props()
        page := view.ProfilePage(view.ProfilePageModel{
            Email: props.EmailAddress.Value(),
            FirstName: props.Name.FirstName(),
            LastName: props.Name.LastName(),
        })
		app := view.Page(view.AppModel{}, page)
		return app.Render(c.Request().Context(), c.Response().Writer)
    }))

    g.PUT("profile", middleware.AuthMiddleware(true)(func(c echo.Context) error {
        userId := c.Get("user_id").(int)

        w := c.Response().Writer
        ctx := c.Request().Context()

        var req postProfileRequest
        if err := c.Bind(&req); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        if err := c.Validate(req); err != nil {
            if errors, ok := err.(validator.ValidationErrors); ok {
                model := view.ProfilePageModel {
                    FirstName: req.FirstName,
                    LastName: req.LastName,
                    Email: req.Email,
                    ValidationErrors: &errors,
                }
                return view.ProfileForm(model).Render(ctx, w)
            } else {
                return err
            }
        }

        user, err := h.userRepo.Get(userId)
        if err != nil {
            return err
        }

        name, err := domain_model.NewUserName(req.FirstName, req.LastName)
        if err != nil {
            return err
        }
        user.ChangeName(name)

        email, err := domain_model.NewUserEmail(req.Email)
        if err != nil {
            return err
        }
        user.ChangeEmail(email)

        if (req.Password != "") {
            user.ChangePassword(req.Password)
        }

        if err := h.userRepo.Commit(*user); err != nil {
            return err
        }

        c.Response().Header().Set("Hx-Location", "/passages")
        c.NoContent(http.StatusNoContent)

        return nil
    }))
}
