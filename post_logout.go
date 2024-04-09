package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func PostLogout(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/logout", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := ctx.SessionStore.Get(r, "session")
		if err != nil {
			return err
		}

		if session != nil {
			delete(session.Values, "user_id")
            if err := session.Save(r, w); err != nil {
                return err
			}
		}

		w.Header().Set("Hx-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)

        return nil
	})).Methods("Post")
}
