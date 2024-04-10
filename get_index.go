package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetIndex(router *mux.Router, ctx *ServerContext) {
	router.Handle("/", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderIndex()
	}))).Methods("Get")
}
