package main

import (
	"net/http"
	"strconv"
	"time"
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
	ID      string
	user_id *int32
}

func GetSession(r *http.Request, ctx *ServerContext) (*Session, error) {
	session, err := ctx.SessionStore.Get(r, "session")
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, nil
	}

	user_id, ok := session.Values["user_id"].(int32)
	if !ok {
		return nil, nil
	}

	return &Session{ID: session.ID, user_id: &user_id}, nil
}

