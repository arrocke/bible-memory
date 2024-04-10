package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages/new", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
		if r.Header.Get("Hx-Current-Url") == "" {
            return engine.RenderCreatePassage(userId, GetClientDate(r))
		} else {
            return engine.RenderCreatePassagePartial()
		}
	}))).Methods("Get")
}
