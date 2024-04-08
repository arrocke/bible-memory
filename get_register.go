package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRegister(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session != nil {
			http.Redirect(w, r, "/passages", http.StatusFound)
			return
		}

        view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderRegister()
	}).Methods("Get")
}
