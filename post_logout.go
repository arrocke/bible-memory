package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func PostLogout(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/logout", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        if _, err := ctx.SessionManager.LogOut(w, r); err != nil {
            return err
        }

		w.Header().Set("Hx-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)

        return nil
	}))).Methods("Post")
}


