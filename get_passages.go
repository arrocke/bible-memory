package main

import (
	"net/http"
	"main/view"

	"github.com/gorilla/mux"
)


func GetPassages(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)
        clientDate := GetClientDate(r)

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
        return engine.RenderPassages(userId, clientDate)
	}))).Methods("GET")
}
