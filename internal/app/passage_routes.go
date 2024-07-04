package app

import (
	"errors"
	"fmt"
	"main/internal/db"
	"main/internal/model"
	"main/internal/view"
	"net/http"
	"regexp"
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
	e.GET("/", func(c echo.Context) error {
        currentUrl := CurrentUrl(c)
        if currentUrl != nil && (strings.HasPrefix(currentUrl.Path, "/passages") || currentUrl.Path == "/") {
            Retarget(c, "#passage-view")
            return c.NoContent(http.StatusOK)
        } else {
            component, err := ctx.CreatePassagesList(c, nil, true)
            if err != nil {
                return err
            }

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
        }
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

        component, err := ctx.CreatePassagesList(c, nil, true)
        if err != nil {
            return err
        }

        PushUrl(c, "/")
        Retarget(c, "#view")
        return RenderComponent(c, component)
	}, AuthMiddleware(true))

	e.GET("/passages/new", func(c echo.Context) error {
        passageView := view.AddPassageView(view.AddPassageViewModel{})

        currentUrl := CurrentUrl(c)
        if currentUrl != nil && (strings.HasPrefix(currentUrl.Path, "/passages") || currentUrl.Path == "/") {
            Retarget(c, "#passage-view")
            return RenderComponent(c, passageView)
        } else {
            component, err := ctx.CreatePassageView(c, passageView)
            if err != nil {
                return err
            }
            return RenderComponent(c, view.Html(component))
        }
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

        passageView := view.EditPassageView(view.EditPassageViewModel{
            Id: passage.Id,
            Reference: passage.Reference.String(),
            Text: passage.Text,
            Interval: passage.Interval,
            NextReview: passage.NextReview,
        })

        currentUrl := CurrentUrl(c)
        if currentUrl != nil && (strings.HasPrefix(currentUrl.Path, "/passages") || currentUrl.Path == "/") {
            Retarget(c, "#passage-view")
            return RenderComponent(c, passageView)
        } else {
            component, err := ctx.CreatePassageView(c, passageView)
            if err != nil {
                return err
            }
            return RenderComponent(c, view.Html(component))
        }
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

        component, err := ctx.CreatePassagesList(c, nil, true)
        if err != nil {
            return err
        }

        PushUrl(c, "/")
        Retarget(c, "#view")
        Reswap(c, "innerHTML")
        return RenderComponent(c, component)
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

	wordRegex := regexp.MustCompile(`(?:(\d+)\s?)?([^A-Za-zÀ-ÖØ-öø-ÿ\s]+)?([A-Za-zÀ-ÖØ-öø-ÿ]+(?:(?:'|’|-)[A-Za-zÀ-ÖØ-öø-ÿ]+)?(?:'|’)?)([^A-Za-zÀ-ÖØ-öø-ÿ0-9]*\s+)?`)

    createReview := func(c echo.Context, passage model.Passage, complete bool) templ.Component {
        wordMatches := wordRegex.FindAllStringSubmatch(passage.Text, -1)
        words := make([]view.ReviewWord, len(wordMatches))
		for i, match := range wordMatches {
			words[i] = view.ReviewWord{
				Number:      match[1],
				Prefix:      match[2],
				Word:        match[3],
				Suffix:      match[4],
			}
		}

        now := GetClientDate(c)

        passageView := view.ReviewPassageView(view.ReviewPassageViewModel{
            Id: passage.Id,
            Reference: passage.Reference.String(),
            Words: words,
            HardInterval: passage.NextInterval(2, now),
            GoodInterval: passage.NextInterval(3, now),
            EasyInterval: passage.NextInterval(4, now),
            AlreadyReviewed: passage.ReviewedAt.Equal(now),
            Complete: complete,
            NextReview: passage.NextReview,
        })

        return passageView

    }

    e.GET("/passages/:id/review", func(c echo.Context) error {
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

        review := createReview(c, passage, false)

        currentUrl := CurrentUrl(c)
        if currentUrl != nil && (strings.HasPrefix(currentUrl.Path, "/passages") || currentUrl.Path == "/") {
            Retarget(c, "#passage-view")
            return RenderComponent(c, review)
        } else {
            component, err := ctx.CreatePassageView(c, review)
            if err != nil {
                return err
            }
            return RenderComponent(c, view.Html(component))
        }
	}, AuthMiddleware(true))

    type postReviewRequest struct {
        Mode string `form:"mode"`
        Grade int `form:"grade"`
    }

    e.POST("/passages/:id/review", func(c echo.Context) error {
        userId, err := GetAuthenticatedUser(c)
        if err != nil {
            return err
        }

        var id int
        if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
            return Redirect(c, "/")
        }

		var req postReviewRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
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

        if req.Mode == "review" {
            passage.Review(req.Grade, GetClientDate(c))

            if err := ctx.PassageRepo.Update(c.Request().Context(), passage); err != nil {
                return err
            }
        }

        review := createReview(c, passage, false)

        component, err := ctx.CreatePassagesList(c, review, true)
        if err != nil {
            return err
        }

        Retarget(c, "#view")
        return RenderComponent(c, component)
    }, AuthMiddleware(true))
}
