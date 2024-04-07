package main

import (
	"context"
	"errors"
	"main/domain_model"
	"main/view"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
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

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, reviewed_at, interval FROM passage WHERE id = $1 AND user_id = $2"
		rows, _ := ctx.Conn.Query(context.Background(), query, id, *session.user_id)
		defer rows.Close()

		passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[PassageModel])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		now := GetClientDate(r)

		model := view.ReviewPassagePageModel{
			Id:              passage.Id,
			Reference:       domain_model.PassageReference{passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse}.String(),
            Text: passage.Text,
			AlreadyReviewed: passage.ReviewedAt != nil && passage.ReviewedAt.Equal(now),
            /*
			HardInterval:    GetNextInterval(now, 2, passage.Interval, passage.ReviewedAt),
			GoodInterval:    GetNextInterval(now, 3, passage.Interval, passage.ReviewedAt),
			EasyInterval:    GetNextInterval(now, 4, passage.Interval, passage.ReviewedAt),
            */
		}

		if r.Header.Get("Hx-Current-Url") == "" {
            model, err := LoadPassagesPageModel(ctx.Conn, *session.user_id, GetClientDate(r), model)
            if err != nil {
                http.Error(w, "Database Error", http.StatusInternalServerError)
                return
            }

            view.App(model).Render(r.Context(), w)
		} else {
            view.ReviewPassagePage(model).Render(r.Context(), w)
		}
	}).Methods("Get")
}
