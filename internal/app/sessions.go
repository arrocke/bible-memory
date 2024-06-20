package app

import (
	"strconv"
	"time"

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

func LogOut(c echo.Context) error {
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

func GetClientDate(c echo.Context) time.Time {
	location := time.FixedZone("Temp", GetClientTZ(c)*60)
	now := time.Now().In(location)

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func GetClientTZ(c echo.Context) int {
	cookieVal, err := c.Cookie("tzOffset")
	if err != nil {
		return 0
	}
	parsedVal, err := strconv.ParseInt(cookieVal.Value, 10, 32)
	if err != nil {
		return 0
	}
	return int(parsedVal)
}
