package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetLogin(router *mux.Router, conn *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/layout.html"))

	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}).Methods("Get")
}
