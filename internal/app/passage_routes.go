package app

import (
	"main/internal/model"
	"main/internal/view"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type postPassagesRequest struct {
	Reference string `form:"reference" validate:"required,reference"`
	Text      string `form:text validate:"required"`
}

func passageRoutes(e *echo.Echo, ctx ServerContext) {
	RenderPassagesView := func(c echo.Context, viewComponent templ.Component) error {
		userId, err := GetAuthenticatedUser(c)
		if err != nil {
			return err
		}

		passages, err := ctx.PassageRepo.GetPassagesForOwner(c.Request().Context(), userId)
		if err != nil {
			return err
		}

		return RenderHtml(c, view.PassagesView(view.PassagesViewModel{
			Passages: passages,
			View:     viewComponent,
		}))
	}

	e.GET("/", func(c echo.Context) error {
		return RenderPassagesView(c, nil)
	}, AuthMiddleware(true))

	e.GET("/passages/new", func(c echo.Context) error {
		return RenderPassagesView(c, view.AddPassageView(view.AddPassageViewModel{}))
	})

	e.POST("/passages", func(c echo.Context) error {
		var req postPassagesRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				model := view.AddPassageViewModel{
                    Text: req.Text,
                    Reference: req.Reference,
					Errors: &errors,
				}
				return RenderComponent(c, view.AddPassageForm(model))
			} else {
				return err
			}
		}

        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        passage, err := model.CreatePassage(req.Reference, req.Text, userId)
        if err != nil {
            return err
        }

        print(passage.Reference.StartChapter)

        if err := ctx.PassageRepo.Create(c.Request().Context(), passage); err != nil {
            return err
        }

        return Redirect(c, "/")
	})
}
