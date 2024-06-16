package app

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionMiddleware(secret string) echo.MiddlewareFunc {
    return session.Middleware(sessions.NewCookieStore([]byte(secret)))
}

func AuthMiddleware(authRequired bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
            userId, err := GetAuthenticatedUser(c)
            if err != nil {
                return err
            }

            if authRequired && userId <= 0 {
                return RedirectWithRefresh(c, "/login")
            } else if !authRequired && userId > 0 {
                return RedirectWithRefresh(c, "/")
            }

			return next(c)
		}
	}
}

func LogIn(c echo.Context, userId int) error {
    sess, err := session.Get("session", c)
    if err != nil {
        return err
    }

    sess.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   86400 * 7,
        HttpOnly: true,
    }
    sess.Values["user_id"] = userId

    if err := sess.Save(c.Request(), c.Response()); err != nil {
        return err
    }

    return nil
}

func LogOut(c echo.Context, userId string) error {
    sess, err := session.Get("session", c)
    if err != nil {
        return err
    }

    sess.Options.MaxAge = -1
    delete(sess.Values, "user_id")

    if err := sess.Save(c.Request(), c.Response()); err != nil {
        return err
    }

    return nil
}

func GetAuthenticatedUser(c echo.Context) (int, error) {
    sess, err := session.Get("session", c)
    if err != nil {
        return 0, err
    }

    if userId, ok := sess.Values["user_id"].(int); ok {
        return userId, nil
    } else {
        return 0, nil
    }
}
