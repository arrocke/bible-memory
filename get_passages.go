package main

import (
	"net/http"
	"main/view"

	"github.com/gorilla/mux"
)


func GetPassages(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
            return err
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

        if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderPassages(int(*session.user_id), GetClientDate(r)); err != nil {
            return err
        }

        return nil
	})).Methods("GET")
}
