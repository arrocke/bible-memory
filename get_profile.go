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

	router.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)

        if err = engine.RenderProfile((int)(*session.user_id)); err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
        }
	}).Methods("Get")
}
