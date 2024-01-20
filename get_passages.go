package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageListItem struct {
	Id        int32
	Reference string
	Level     int32
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
	}

	query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, level FROM passage ORDER BY id ASC"
	rows, _ := conn.Query(context.Background(), query)
	defer rows.Close()

	passages, err := pgx.CollectRows(rows, pgx.RowToStructByName[PassageModel])
	if err != nil {
		return nil, err
	}

	listItems := make([]PassageListItem, len(passages))
	for i, passage := range passages {
		listItems[i] = PassageListItem{
			Id:        passage.Id,
			Level:     passage.Level,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
		}
	}

	return listItems, nil
}

func GetPassages(router *mux.Router, conn *pgxpool.Pool) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Level        int32
	}

	type Passage struct {
		Id        int32
		Reference string
		Level     int32
	}

	type TemplateData struct {
		Passages []PassageListItem
	}

	tmpl := template.Must(template.ParseFiles("templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		passages, err := LoadPassageList(conn)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", TemplateData{Passages: passages})
	}).Methods("GET")
}
