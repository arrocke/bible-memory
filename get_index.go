package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetIndex(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/public_index.html", "templates/layout.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}

		if session == nil {
			tmpl.ExecuteTemplate(w, "layout.html", nil)
		} else {
			http.Redirect(w, r, "/passages", http.StatusFound)
		}
	}).Methods("Get")
}
