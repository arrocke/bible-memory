package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogin(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			println(err.Error())
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}

		if session == nil {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLogin()
		} else {
			http.Redirect(w, r, "/passages", http.StatusFound)
		}
	}).Methods("Get")
}
