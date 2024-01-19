package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPassages(router *mux.Router, conn *pgxpool.Pool) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
	}

	type Passage struct {
		Id        int32
		Reference string
		Level     int32
	}

	type TemplateData struct {
		Passages []Passage
	}

	tmpl := template.Must(template.ParseFiles("templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse FROM public.passage ORDER BY id ASC"
		rows, _ := conn.Query(context.Background(), query)
		defer rows.Close()

		passages, err := pgx.CollectRows(rows, pgx.RowToStructByName[PassageModel])
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		templateData := TemplateData{Passages: make([]Passage, len(passages))}
		for i, passage := range passages {
			templateData.Passages[i] = Passage{
				Id:        passage.Id,
				Level:     1,
				Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
			}
		}

		tmpl.ExecuteTemplate(w, "layout.html", templateData)
	})
}
