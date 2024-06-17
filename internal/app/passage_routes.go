package app

import (
	"main/internal/model"
	"main/internal/view"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func passageRoutes(e *echo.Echo, ctx ServerContext) {
    GetPassageId := func(c echo.Context) int {
        i64, err := strconv.ParseInt(c.Param("id"), 10, 32)
        if err != nil {
            return 0
        } else {
            return int(i64)
        }
    }

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

    type postPassagesRequest struct {
        Reference string `form:"reference" validate:"required,reference"`
        Text      string `form:"text" validate:"required"`
    }

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

        if err := ctx.PassageRepo.Create(c.Request().Context(), passage); err != nil {
            return err
        }

        return Redirect(c, "/")
	})

	e.GET("/passages/new", func(c echo.Context) error {
		return RenderPassagesView(c, view.AddPassageView(view.AddPassageViewModel{}))
	})

    e.GET("/passages/:id", func (c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        id := GetPassageId(c)
        if id == 0 {
            return Redirect(c, "/")
        }

        passage, err := ctx.PassageRepo.GetPassageById(c.Request().Context(), id)
        if err != nil {
            return err
        }

        if passage.Owner != userId {
            return Redirect(c, "/")
        }

		return RenderPassagesView(c, view.EditPassageView(view.EditPassageViewModel{
            Passage: passage,
        }))
    })
}
