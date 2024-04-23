package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type SessionManager struct {
	store *sessions.CookieStore
}

func InitSessions(secret string) SessionManager {
	return SessionManager{
		store: sessions.NewCookieStore([]byte(secret)),
	}
}

func (m SessionManager) Middleware() echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
        return func (c echo.Context) error {
            session, err := m.store.Get(c.Request(), "session")
            if err != nil {
                return err
            }

            if userId, ok := session.Values["user_id"].(int); ok {
                c.Set("user_id", userId)
            }
            
            return next(c)
        }
    }
}

func (m SessionManager) LogIn(c echo.Context, userId int) error {
	sess, err := m.store.New(c.Request(), "session")
    if err != nil {
        return err
    }

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
    sess.Values["user_id"] = userId

    return sess.Save(c.Request(), c.Response())
}

func (m SessionManager) LogOut(c echo.Context) error {
	sess, err := m.store.New(c.Request(), "session")
    if err != nil {
        return err
    }

    return sess.Save(c.Request(), c.Response())
}

func GetUserId(c echo.Context) int {
    id := c.Get("user_id")
    if id == nil {
        return 0
    } else {
        return id.(int)
    }
}
