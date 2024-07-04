package app

import (
	"net/http"
	"net/url"

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

func Retarget(c echo.Context, target string) {
    c.Response().Header().Set("hx-retarget", target)
}

func IsHtmxRequest(c echo.Context) bool {
    return c.Request().Header.Get("hx-request") == "true"
}

func CurrentUrl(c echo.Context) *url.URL {
    header := c.Request().Header.Get("hx-current-url")
    if header == "" {
        return nil
    } else {
        parsedUrl, err := url.Parse(header)
        if err == nil {
            return parsedUrl
        } else {
            return nil
        }
    }
}
