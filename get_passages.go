package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPassages(router *mux.Router) {
	type Passage struct {
		Id        string
		Reference string
		Level     int
	}

	type TemplateData struct {
		Passages []Passage
	}

	tmpl := template.Must(template.ParseFiles("templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
			Passages: []Passage{
				{Id: "1", Reference: "Genesis 1:1-5", Level: 1},
				{Id: "2", Reference: "Genesis 1:5-10", Level: 2},
			},
		})
	})
}
