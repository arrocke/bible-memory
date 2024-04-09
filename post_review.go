package main

import (
	"main/services"
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}/review", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

		if r.FormValue("mode") != "review" {
            if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(int(*session.user_id), GetClientDate(r)); err != nil {
                return err
            }
			return nil
		}

		grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return nil
		}

		tz := GetClientTZ(r)

		if err := ctx.PassageService.Review(services.ReviewPassageRequest{
			Id:    int(id),
			Grade: int(grade),
			Tz:    tz,
		}); err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return nil
		}

        if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(int(*session.user_id), GetClientDate(r)); err != nil {
            return err
        }

        return nil
	})).Methods("Post")
}
