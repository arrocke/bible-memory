package app

import (
	"errors"
	"main/internal/db"
	"main/internal/model"
	"main/internal/view"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type date time.Time

func (d *date) UnmarshalParam(param string) error {
    parsed, err := time.Parse("2006-01-02", param)
    if err != nil {
        *d = date(time.Time{})
    }
    *d = date(parsed)
    return nil
}

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
            Id: passage.Id,
            Reference: passage.Reference.String(),
            Text: passage.Text,
            Interval: passage.Interval,
            NextReview: passage.NextReview,
        }))
    })

    type putPassageRequest struct {
        Id int `param:"id"`
        Reference string `form:"reference" validate:"required,reference"`
        ReviewAt date `form:"review_at"`
        Interval int `form:"interval" validate:"omitempty,min=1"`
        Text string `form:"text" validate:"required"`
    }

    e.PUT("/passages/:id", func (c echo.Context) error {
		var req putPassageRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				model := view.EditPassageViewModel{
                    Id: req.Id,
                    Text: req.Text,
                    Reference: req.Reference,
                    Interval: req.Interval,
                    NextReview: (time.Time)(req.ReviewAt),
					Errors: &errors,
				}
				return RenderComponent(c, view.EditPassageForm(model))
			} else {
				return err
			}
		}

        passage, err := ctx.PassageRepo.GetPassageById(c.Request().Context(), req.Id)
        if err != nil {
            if errors.Is(err, db.NotFoundError) {
                return Redirect(c, "/")
            } else {
                return err
            }
        }

        // We don't have to handle this error because the request validation should prevent this from returning an error
        reference, err := model.ParseReference(req.Reference)
        if err != nil {
            return err
        }

        passage.Reference = reference
        passage.Text = req.Text
        passage.Interval = req.Interval
        passage.NextReview = time.Time(req.ReviewAt)

        if err := ctx.PassageRepo.Update(c.Request().Context(), passage); err != nil {
            return err
        }

        return Redirect(c, "/")
    })
}
