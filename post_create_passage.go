package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PostCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		reference, err := ParseReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}

		var id int32
		query := "INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse, text) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
		err = ctx.Conn.QueryRow(context.Background(), query, reference.Book, reference.StartChapter, reference.StartVerse, reference.EndChapter, reference.EndVerse, r.FormValue("text")).Scan(&id)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
