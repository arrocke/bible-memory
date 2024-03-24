package main

import (
	"fmt"
	"main/domain_model"
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

		reference, err := domain_model.ParsePassageReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}

		passage := domain_model.NewPassage(
            reference,
            r.FormValue("text"),
            int(*session.user_id),
        )

		if err = ctx.PassageRepo.Commit(&passage); err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", passage.Id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
