package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRegister(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/register", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}
		if session != nil {
			http.Redirect(w, r, "/passages", http.StatusFound)
			return nil
		}

        if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderRegister(); err != nil {
            return err
        }

        return nil
	})).Methods("Get")
}
