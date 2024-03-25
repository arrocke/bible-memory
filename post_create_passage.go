package main

import (
	"fmt"
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func PostCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

        id, err := ctx.PassageService.Create(services.CreatePassageRequest{
            Reference: r.FormValue("reference"),
            Text: r.FormValue("text"),
            UserId: int(*session.user_id),
        })
        if err != nil {
            http.Error(w, "Error", http.StatusBadRequest)
        }

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
