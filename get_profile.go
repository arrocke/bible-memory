package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProfile(router *mux.Router, ctx *ServerContext) {
	type userModel struct {
		Email     string
		FirstName string
		LastName  string
	}

	router.HandleFunc("/users/profile", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
        if err = engine.RenderProfile((int)(*session.user_id)); err != nil {
            return err
        }

        return nil
	})).Methods("Get")
}
