package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetIndex(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}

		if session == nil {
			view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderIndex()
		} else {
			http.Redirect(w, r, "/passages", http.StatusFound)
		}

        return nil
	})).Methods("Get")
}
