package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRegister(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/register.html", "templates/layout.html"))

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session != nil {
			http.Redirect(w, r, "/passages", http.StatusFound)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", LayoutTemplateData{})
	}).Methods("Get")
}
