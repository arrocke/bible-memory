package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginTemplateData struct {
	Error string
	Email string
	LayoutTemplateData
}

var LoginTemplate = template.Must(template.ParseFiles("templates/login.html", "templates/layout.html"))

func GetLogin(router *mux.Router, ctx *ServerContext) {
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

		LoginTemplate.ExecuteTemplate(w, "layout.html", LoginTemplateData{})
	}).Methods("Get")
}
