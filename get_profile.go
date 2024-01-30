package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProfile(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/profile.html", "templates/layout.html"))

	type UserProfile struct {
		Email     string
		FirstName string
		LastName  string
	}

	type TemplateData struct {
		LayoutTemplateData
		User UserProfile
	}

	router.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		templateData, err := LoadLayoutTemplateData(ctx.Conn, session.user_id)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", templateData)
	}).Methods("Get")
}
