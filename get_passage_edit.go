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
)

func GetPassageEdit(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Text         string
		ReviewAt     *time.Time
		Interval     *int
	}

	type PartialTemplateData struct {
		Id        int32
		Reference string
		Text      string
		ReviewAt  string
		Interval  *int
	}

	type TemplateData struct {
		PartialTemplateData
		PassagesTemplateData
	}

	tmpl := template.Must(template.ParseFiles("templates/edit_passage_partial.html", "templates/edit_passage.html", "templates/passage_list_partial.html", "templates/passages.html", "templates/layout.html"))

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

		partialTemplateData := PartialTemplateData{
			Id:        passage.Id,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
			Text:      passage.Text,
			Interval:  passage.Interval,
			ReviewAt:  reviewAt,
		}

		if r.Header.Get("Hx-Current-Url") == "" {
			passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
				partialTemplateData,
				*passagesTemplateData,
			})
		} else {
			tmpl.ExecuteTemplate(w, "edit_passage_partial.html", partialTemplateData)
		}
	}).Methods("Get")
}
