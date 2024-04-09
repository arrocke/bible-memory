package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogin(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/login", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}

		if session == nil {
            if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLogin(); err != nil {
                return err
            }
		} else {
			http.Redirect(w, r, "/passages", http.StatusFound)
		}

        return nil
	})).Methods("Get")
}
