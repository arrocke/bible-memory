package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func PostLogout(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/users/logout", func(w http.ResponseWriter, r *http.Request) {
		session, err := ctx.SessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		if session != nil {
			delete(session.Values, "user_id")
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, "Session Error", http.StatusInternalServerError)
			}
		}

		w.Header().Set("Hx-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
