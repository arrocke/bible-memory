package app

import (
	"net/http"

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
