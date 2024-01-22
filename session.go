package main

import (
	"net/http"
)

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
