package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPassageReview(router *mux.Router) {
	type TemplateData struct {
		Id string
	}

	tmpl := template.Must(template.ParseFiles("templates/get_review.html", "templates/layout.html"))

	router.HandleFunc("/passages/{Id}/review", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
			Id: vars["Id"],
		})
	})
}
