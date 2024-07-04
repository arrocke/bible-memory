package app

import (
	"fmt"
	"main/internal/view"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (ctx ServerContext) CreateView(c echo.Context, component templ.Component) (templ.Component, error) {
    userId, err := GetAuthenticatedUser(c)
    if err != nil {
        return nil, err
    }

    if userId == 0 {
        return view.Layout(view.LayoutModel{
            Authenticated: false,
            View: component,
        }), nil
    } else {
        user, err := ctx.UserRepo.GetUserById(c.Request().Context(), userId)
        if err != nil {
            return nil, err
        }

        return view.Layout(view.LayoutModel{
            Authenticated: true,
            User: view.UserModel {
                FirstName: user.FirstName,
                LastName: user.LastName,
            },
            View: component,
        }), nil
    }
}

func (ctx ServerContext) CreatePassagesList(c echo.Context, component templ.Component, menuOpen bool) (templ.Component, error) {
    userId, err := GetAuthenticatedUser(c)
    if err != nil {
        return nil, err
    }

    if userId == 0 {
        return nil, fmt.Errorf("Cannot create passage list with no authenticated user.")
    }

    passages, err := ctx.PassageRepo.GetPassagesForOwner(c.Request().Context(), userId)
    if err != nil {
        return nil, err
    }

    return view.PassagesView(view.PassagesViewModel{
        Passages: passages,
        Now: GetClientDate(c),
        View:     component,
        StartOpen: menuOpen,
    }), nil
}

func (ctx ServerContext) CreatePassageView(c echo.Context, component templ.Component) (templ.Component, error) {
    component, err := ctx.CreatePassagesList(c, component, true)
    if err != nil {
        return nil, err
    }

    return ctx.CreateView(c, component)
}
