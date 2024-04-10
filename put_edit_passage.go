package main

import (
	"fmt"
	"main/services"
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
	router.Handle("/passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

        form, err := parseEditPassageForm(r)
        if err != nil {
            http.Error(w, "Invalid body", http.StatusBadRequest)
            return nil
        }

        if err := ctx.PassageService.Update(services.UpdatePassageRequest{
            Id: int(id),
            Reference: form.Reference,
            Text: form.Text,
            Interval: form.Interval,
            ReviewAt: form.ReviewAt,
        }); err != nil {
            http.Error(w, "Error", http.StatusBadRequest)
            return nil
        }

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)

        return nil
	}))).Methods("Put")
}
