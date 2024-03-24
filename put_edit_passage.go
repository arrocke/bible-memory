package main

import (
	"fmt"
	"main/domain_model"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type editPassageForm struct {
    Reference string;
    Text string;
    Interval *int;
    ReviewAt *time.Time
}

func parseEditPassageForm(r *http.Request) (editPassageForm, error) {
    var form editPassageForm

    if err := r.ParseForm(); err != nil {
        return form, nil
    }

    err := decoder.Decode(&form, r.PostForm)
    return form, err
}

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

        form, err := parseEditPassageForm(r)
        if err != nil {
            http.Error(w, "Invalid body", http.StatusBadRequest)
            return
        }

		reference, err := domain_model.ParsePassageReference(form.Reference)
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}
        passage.SetReference(reference)

        passage.SetText(form.Text)

        if (form.Interval != nil && form.ReviewAt != nil) {
            interval, err := domain_model.NewReviewInterval(*form.Interval)
            if err != nil {
				http.Error(w, "Invalid interval", http.StatusBadRequest)
                return
            }

            nextReview := domain_model.NewReviewTimestamp(*form.ReviewAt)

            passage.SetReviewState(interval, nextReview)
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
