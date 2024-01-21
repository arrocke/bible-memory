package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPassageEdit(router *mux.Router, conn *pgxpool.Pool) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Text         string
		Level        int32
		ReviewAt     *time.Time
	}

	type PartialTemplateData struct {
		Id        int32
		Reference string
		Text      string
		Level     int32
		ReviewAt  string
	}

	type FullTemplateData struct {
		Id        int32
		Reference string
		Text      string
		Level     int32
		ReviewAt  string
		Passages  []PassageListItem
	}

	tmpl := template.Must(template.ParseFiles("templates/edit_passage_partial.html", "templates/edit_passage.html", "templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, level, review_at FROM passage WHERE id = $1"
		rows, _ := conn.Query(context.Background(), query, id)
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

		reviewAt := ""
		if passage.ReviewAt != nil {
			reviewAt = passage.ReviewAt.Format("2006-01-02")
		}

		if r.Header.Get("Hx-Current-Url") == "" {
			passageList, err := LoadPassageList(conn)
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			tmpl.ExecuteTemplate(w, "layout.html", FullTemplateData{
				Id:        passage.Id,
				Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
				Level:     passage.Level,
				ReviewAt:  reviewAt,
				Text:      passage.Text,
				Passages:  passageList,
			})
		} else {
			tmpl.ExecuteTemplate(w, "edit_passage_partial.html", PartialTemplateData{
				Id:        passage.Id,
				Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
				Level:     passage.Level,
				Text:      passage.Text,
				ReviewAt:  reviewAt,
			})
		}

	}).Methods("Get")
}
