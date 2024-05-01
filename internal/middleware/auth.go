package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(requireAuth bool) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
        return func (c echo.Context) error {
            userId := GetUserId(c)

            var redirect string
            if requireAuth && userId == 0 {
                redirect = "/login"
            } else if !requireAuth && userId != 0 {
                redirect = "/passages"
            }

            if redirect != "" {
                hxRequestHeader := c.Request().Header["Hx-Request"]
                if len(hxRequestHeader) > 0 && hxRequestHeader[0] == "true" {
                    c.Response().Header().Set("Hx-Location", redirect)
                    c.NoContent(http.StatusNoContent)
                } else {
                    c.Redirect(http.StatusFound, redirect)
                }

                return nil
            }

            return next(c)
        }
    }
}
