package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCreatePassage(router *mux.Router, conn *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("templates/add_passage.html", "templates/layout.html"))

	router.HandleFunc("/passages/new", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}).Methods("Get")
}
