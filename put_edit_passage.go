package main

import (
	"fmt"
	"main/domain_model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PutEditPassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

        passage, err := ctx.PassageRepo.Get(uint(id))
        if err != nil {
            fmt.Println(err.Error())
        }

		reference, err := domain_model.ParsePassageReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}
        passage.SetReference(reference)

        passage.SetText(r.FormValue("text"))

		intervalStr := r.FormValue("interval")
		reviewAtStr := r.FormValue("review_at")
        if (intervalStr != "" && reviewAtStr != "") {
            interval, err := domain_model.ParseReviewInterval(intervalStr)
            if err != nil {
				http.Error(w, "Invalid interval", http.StatusBadRequest)
                return
            }
            nextReview, err := domain_model.ParseReviewTimestamp(reviewAtStr, "2006-01-02")
            if err != nil {
				http.Error(w, "Invalid review date", http.StatusBadRequest)
                return
            }

            passage.OverrideReviewState(interval, nextReview)
        }

        err = ctx.PassageRepo.Commit(&passage)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Put")
}
