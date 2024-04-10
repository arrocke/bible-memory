package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

func GetClientDate(r *http.Request) time.Time {
	location := time.FixedZone("Temp", GetClientTZ(r)*60)
	now := time.Now().In(location)

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func GetClientTZ(r *http.Request) int {
	cookieVal, err := r.Cookie("tzOffset")
	if err != nil {
		return 0
	}
	parsedVal, err := strconv.ParseInt(cookieVal.Value, 10, 32)
	if err != nil {
		return 0
	}
	return int(parsedVal)
}

type Session struct {
	Id      string
	UserId  *int
}

type SessionManager struct {
    store *sessions.CookieStore
}

func CreateSessionManager(store *sessions.CookieStore) *SessionManager {
    return &SessionManager{store}
}

func (manager SessionManager) LogOut(w http.ResponseWriter, r *http.Request) (Session, error) {
    session, err := manager.store.Get(r, "session")
    if err != nil {
        return Session{}, err
    }

    delete(session.Values, "user_id")
    if err := session.Save(r, w); err != nil {
        return Session{}, err
    }

    return Session{ Id: session.ID }, nil
}

func (manager SessionManager) LogIn(w http.ResponseWriter, r *http.Request, userId int) (Session, error) {
    session, err := manager.store.New(r, "session")
    if err != nil {
        return Session{}, err
    }

    session.Values["user_id"] = userId
    if err := session.Save(r, w); err != nil {
        return Session{}, err
    }

    return Session{
        Id: session.ID,
        UserId: &userId,
    }, nil
}

func (manager SessionManager) SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        session, err := manager.store.Get(r, "session")
        if err != nil {
            http.Error(w, "Error loading session", http.StatusInternalServerError)
            return
        }

        sessionData := Session {
            Id: session.ID,
        }

        userId, ok := session.Values["user_id"].(int)
        if ok {
            sessionData.UserId = &userId
        }

        newRequest := r.WithContext(context.WithValue(r.Context(), "session", sessionData))

        next.ServeHTTP(w, newRequest)
    })
}

func GetSession(r *http.Request) Session {
    return r.Context().Value("session").(Session)
}

func GetUserId(r *http.Request) int {
    session := GetSession(r)
    return *session.UserId
}

func AuthMiddleware(requireAuth bool, next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        session := GetSession(r)

        if requireAuth && session.UserId == nil {
            w.Header().Set("Hx-Redirect", "/")
            w.WriteHeader(http.StatusNoContent)
            return
        } else if !requireAuth && session.UserId != nil {
            w.Header().Set("Hx-Redirect", "/passages")
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

