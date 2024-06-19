package app

import (
	"errors"
	"fmt"
	"main/internal/db"
	"main/internal/model"
	"main/internal/view"
	"net/http"
	"strings"
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
	}, AuthMiddleware(true))

	e.GET("/passages/new", func(c echo.Context) error {
		return RenderPassagesView(c, view.AddPassageView(view.AddPassageViewModel{}))
	}, AuthMiddleware(true))

    e.GET("/passages/:id", func (c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        var id int
        if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
            return Redirect(c, "/")
        }

        passage, err := ctx.PassageRepo.GetPassageById(c.Request().Context(), id)
        if err != nil {
            if errors.Is(err, db.NotFoundError) {
                return Redirect(c, "/")
            } else {
                return err
            }
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
    }, AuthMiddleware(true))

    type putPassageRequest struct {
        Id int `param:"id"`
        Reference string `form:"reference" validate:"required,reference"`
        ReviewAt date `form:"review_at"`
        Interval int `form:"interval" validate:"omitempty,min=1"`
        Text string `form:"text" validate:"required"`
    }

    e.PUT("/passages/:id", func (c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

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

        if passage.Owner != userId {
            return Redirect(c, "/")
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
    }, AuthMiddleware(true))

    e.DELETE("/passages/:id", func (c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        var id int
        if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
            return Redirect(c, "/")
        }

        passage, err := ctx.PassageRepo.GetPassageById(c.Request().Context(), id)
        if err != nil {
            if errors.Is(err, db.NotFoundError) {
                return Redirect(c, "/")
            } else {
                return err
            }
        }

        if passage.Owner != userId {
            return Redirect(c, "/")
        }

        if err := ctx.PassageRepo.Delete(c.Request().Context(), passage.Id); err != nil {
            return err
        }

        currentUrl := c.Request().Header.Get("hx-current-url")
        if strings.HasSuffix(currentUrl, fmt.Sprintf("/passages/%v", passage.Id)) ||
            strings.Contains(currentUrl, fmt.Sprintf("/passages/%v/", passage.Id)) {
            return Redirect(c, "/")
        } else {
            return c.NoContent(http.StatusOK)
        }

    }, AuthMiddleware(true))
}
