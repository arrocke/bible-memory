package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCreatePassage(router *mux.Router, conn *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("templates/add_passage.html", "templates/add_passage_partial.html", "templates/passages.html", "templates/layout.html"))

	type TemplateData struct {
		Passages []PassageListItem
	}

	router.HandleFunc("/passages/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Current-Url") == "" {
			passageList, err := LoadPassageList(conn)
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}

			tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
				Passages: passageList,
			})
		} else {
			tmpl.ExecuteTemplate(w, "add_passage_partial.html", nil)
		}

	}).Methods("Get")
}
