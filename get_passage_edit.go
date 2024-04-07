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

func GetPassageEdit(router *mux.Router, ctx *ServerContext) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		Text         string
		ReviewAt     *time.Time
		Interval     *int
	}

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

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, review_at, interval FROM passage WHERE id = $1 AND user_id = $2"
		rows, _ := ctx.Conn.Query(context.Background(), query, id, *session.user_id)
		defer rows.Close()

		passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[passageModel])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		model := view.EditPassagePageModel {
			Id:        passage.Id,
			Reference: domain_model.PassageReference{passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse}.String(),
			Text:      passage.Text,
			Interval:  passage.Interval,
			ReviewAt:  passage.ReviewAt,
		}

		if r.Header.Get("Hx-Current-Url") == "" {
            model, err := LoadPassagesPageModel(ctx.Conn, *session.user_id, GetClientDate(r), model)
            if err != nil {
                http.Error(w, "Database Error", http.StatusInternalServerError)
                return
            }

            view.App(model).Render(r.Context(), w)
		} else {
            view.EditPassagePage(model).Render(r.Context(), w)
		}
	}).Methods("Get")
}
