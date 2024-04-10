package main

import (
	"fmt"
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func PostCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

        id, err := ctx.PassageService.Create(services.CreatePassageRequest{
            Reference: r.FormValue("reference"),
            Text: r.FormValue("text"),
            UserId: userId,
        })
        if err != nil {
            http.Error(w, "Error", http.StatusBadRequest)
        }

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)

        return nil
	}))).Methods("Post")
}
