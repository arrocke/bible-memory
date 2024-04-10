package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRegister(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
        return engine.RenderRegister(r.Header.Get("Hx-Request") != "true")
	}))).Methods("Get")
}
