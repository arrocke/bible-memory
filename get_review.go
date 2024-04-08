package main

import (
	"main/view"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetPassageReview(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		Text         string
		ReviewedAt   *time.Time
		Interval     *int
	}

	router.HandleFunc("/passages/{Id}/{Mode}", func(w http.ResponseWriter, r *http.Request) {
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

		if r.Header.Get("Hx-Current-Url") == "" {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewPassage(int(*session.user_id), int(id), GetClientDate(r))
		} else {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewPassagePartial(int(*session.user_id), int(id), GetClientDate(r))
		}
	}).Methods("Get")
}
