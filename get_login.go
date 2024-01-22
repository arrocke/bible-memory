package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogin(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/layout.html"))

	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			println(err.Error())
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
