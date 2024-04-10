package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogin(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/login", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
        return engine.RenderLogin(r.Header.Get("Hx-Request") != "true")
	}))).Methods("Get")
}
