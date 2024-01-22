package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCreatePassage(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/add_passage.html", "templates/add_passage_partial.html", "templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages/new", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if r.Header.Get("Hx-Current-Url") == "" {
			templateData, err := LoadPassagesTemplateData(ctx.Conn)
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}

			tmpl.ExecuteTemplate(w, "layout.html", templateData)
		} else {
			tmpl.ExecuteTemplate(w, "add_passage_partial.html", nil)
		}

	}).Methods("Get")
}
