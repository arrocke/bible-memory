package main

import (
	"main/services"
	"main/view"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var REVIEW_MAP = [...]int{1, 1, 1, 2, 2, 3, 5, 8, 13, 21, 34, 55}

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

	router.HandleFunc("/passages/{Id}/review", func(w http.ResponseWriter, r *http.Request) {
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

		if r.FormValue("mode") != "review" {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(int(*session.user_id), GetClientDate(r))
			return
		}

		grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return
		}

		tz := GetClientTZ(r)

		if err := ctx.PassageService.Review(services.ReviewPassageRequest{
			Id:    int(id),
			Grade: int(grade),
			Tz:    tz,
		}); err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

        view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(int(*session.user_id), GetClientDate(r))
	}).Methods("Post")
}
