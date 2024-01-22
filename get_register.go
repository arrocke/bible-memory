package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRegister(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/register.html", "templates/layout.html"))

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}).Methods("Get")
}
