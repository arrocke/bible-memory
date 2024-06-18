package app

import (
	"main/internal/model"
	"main/internal/view"
	"net/http"
	"strconv"
	"time"

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
            Id: passage.Id,
            Reference: passage.Reference.String(),
            Text: passage.Text,
            Interval: passage.Interval,
            NextReview: passage.NextReview,
        }))
    })

    type putPassageRequest struct {
        Id int
        Reference string `validate:"required,reference"`
        ReviewAt *time.Time
        Interval *int `validate:"omitnil,min=1"`
        Text string `validate:"required"`
    }

    e.PUT("/passages/:id", func (c echo.Context) error {
		var req putPassageRequest

        binder := echo.FormFieldBinder(c)
        errs := binder.String("reference", &req.Reference).
            String("text", &req.Text).
            CustomFunc("interval", func(values []string) []error {
                if len(values) == 0 || values[0] == "" {
                    return nil
                } else {
                    interval64, err := strconv.ParseInt(values[0], 10, 32)
                    if err != nil {
                        println(err.Error())
                        return []error{echo.NewBindingError("interval", values[0:1], "failed to decode int", err)}
                    }
                    interval := int(interval64)
                    req.Interval = &interval
                    return nil
                }
            }).
            CustomFunc("review_at", func(values []string) []error {
                if len(values) == 0  || values[0] == ""{
                    return nil
                } else {
                    reviewAt, err := time.Parse("2006-01-02", values[0])
                    if err != nil {
                        println(err.Error())
                        return []error{echo.NewBindingError("review_at", values[0:1], "failed to decode date", err)}
                    }
                    req.ReviewAt = &reviewAt
                    return nil
                }
            }).
            BindErrors()
		if len(errs) > 0 {
			return c.String(http.StatusBadRequest, "bad request")
		}

        binder = echo.PathParamsBinder(c)
        errs = binder.Int("id", &req.Id).BindErrors()
		if len(errs) > 0 {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				model := view.EditPassageViewModel{
                    Id: req.Id,
                    Text: req.Text,
                    Reference: req.Reference,
                    Interval: req.Interval,
                    NextReview: req.ReviewAt,
					Errors: &errors,
				}
				return RenderComponent(c, view.EditPassageForm(model))
			} else {
				return err
			}
		}

				model := view.EditPassageViewModel{
                    Id: req.Id,
                    Text: req.Text,
                    Reference: req.Reference,
                    Interval: req.Interval,
                    NextReview: req.ReviewAt,
				}
				return RenderComponent(c, view.EditPassageForm(model))

        return nil
    })
}
