package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageListItem struct {
	Id        int32
	Reference string
	Level     int32
	ReviewAt  string
}

func LoadPassageList(conn *pgxpool.Pool) ([]PassageListItem, error) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Level        int32
		ReviewAt     *time.Time
	}

	query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, level, review_at FROM passage ORDER BY id ASC"
	rows, _ := conn.Query(context.Background(), query)
	defer rows.Close()

	passages, err := pgx.CollectRows(rows, pgx.RowToStructByName[PassageModel])
	if err != nil {
		println(err.Error())
		return nil, err
	}

	listItems := make([]PassageListItem, len(passages))
	for i, passage := range passages {
		reviewAt := ""
		if passage.ReviewAt != nil {
			reviewAt = passage.ReviewAt.Format("01-02-2006")
		}
		listItems[i] = PassageListItem{
			Id:        passage.Id,
			Level:     passage.Level,
			ReviewAt:  reviewAt,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
		}
	}

	return listItems, nil
}

func GetPassages(router *mux.Router, ctx *ServerContext) {
	type TemplateData struct {
		Passages []PassageListItem
	}

	tmpl := template.Must(template.ParseFiles("templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		session, err := ctx.SessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}

		user_id := session.Values["user_id"]
		fmt.Println(user_id)

		passages, err := LoadPassageList(ctx.Conn)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", TemplateData{Passages: passages})
	}).Methods("GET")
}
