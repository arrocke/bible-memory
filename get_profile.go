package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProfile(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/profile", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
        return engine.RenderProfile(userId)
	}))).Methods("Get")
}
