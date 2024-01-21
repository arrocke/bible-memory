package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRegister(router *mux.Router, conn *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("templates/register.html", "templates/layout.html"))

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}).Methods("Get")
}
