package app

import (
	"main/internal/view"

	"github.com/labstack/echo/v4"
)

func passageRoutes(e *echo.Echo, ctx ServerContext) {
	e.GET("/", func(c echo.Context) error {
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
        }, nil))
	}, AuthMiddleware(true))
}
