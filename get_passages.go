package main

import (
	"context"
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
type PassagesTemplateData struct {
	Passages []PassageListItem
	LayoutTemplateData
}

func LoadPassagesTemplateData(conn *pgxpool.Pool, layoutTemplateData LayoutTemplateData) (*PassagesTemplateData, error) {
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

	templateData := PassagesTemplateData{
		Passages:           make([]PassageListItem, len(passages)),
		LayoutTemplateData: layoutTemplateData,
	}
	for i, passage := range passages {
		reviewAt := ""
		if passage.ReviewAt != nil {
			reviewAt = passage.ReviewAt.Format("01-02-2006")
		}
		templateData.Passages[i] = PassageListItem{
			Id:        passage.Id,
			Level:     passage.Level,
			ReviewAt:  reviewAt,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
		}
	}

	return &templateData, nil
}

func GetPassages(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		data, err := LoadPassagesTemplateData(ctx.Conn, LayoutTemplateData{IsLoggedIn: true})
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", data)
	}).Methods("GET")
}
