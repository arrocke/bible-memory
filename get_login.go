package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogin(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/layout.html"))

	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}).Methods("Get")
}
