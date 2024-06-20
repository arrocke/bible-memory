package app

import (
	"main/internal/view"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Redirect(c echo.Context, path string) error {
    if c.Request().Header.Get("hx-request") == "true" {
        c.Response().Header().Set("hx-location", path)
        return c.NoContent(http.StatusOK)
    } else {
        c.Response().Header().Set("Location", path)
        return c.NoContent(http.StatusSeeOther)
    }
}

func RedirectWithRefresh(c echo.Context, path string) error {
    if c.Request().Header.Get("hx-request") == "true" {
        c.Response().Header().Set("hx-redirect", path)
        return c.NoContent(http.StatusOK)
    } else {
        c.Response().Header().Set("Location", path)
        return c.NoContent(http.StatusSeeOther)
    }
}

func RenderComponent(c echo.Context, component templ.Component) error {
    return component.Render(c.Request().Context(), c.Response().Writer)
}

func RenderHtml(c echo.Context, component templ.Component) error {
    return RenderComponent(c, view.Html(component))
}

func (ctx ServerContext) RenderView(c echo.Context, viewComponent templ.Component) error {
    userId, err := GetAuthenticatedUser(c)
    if err != nil {
        return err
    }

    if userId == 0 {
        return RenderHtml(c, view.Layout(view.LayoutModel{
            Authenticated: false,
            View: viewComponent,
        }))
    } else {
        user, err := ctx.UserRepo.GetUserById(c.Request().Context(), userId)
        if err != nil {
            return err
        }

        return RenderHtml(c, view.Layout(view.LayoutModel{
            Authenticated: true,
            User: view.UserModel {
                FirstName: user.FirstName,
                LastName: user.LastName,
            },
            View: viewComponent,
        }))
    }
}

func (ctx ServerContext) RenderPassagesView(c echo.Context, viewComponent templ.Component) error {
    userId, err := GetAuthenticatedUser(c)
    if err != nil {
        return err
    }

    passages, err := ctx.PassageRepo.GetPassagesForOwner(c.Request().Context(), userId)
    if err != nil {
        return err
    }

    startOpen := false
    if viewComponent == nil {
        startOpen = true
    }

    return ctx.RenderView(c, view.PassagesView(view.PassagesViewModel{
        Passages: passages,
        Now: GetClientDate(c),
        View:     viewComponent,
        StartOpen: startOpen,
    }))
}

