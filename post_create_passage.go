package main

import (
	"context"
	"fmt"
	"main/domain_model"
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

		reference, err := domain_model.ParsePassageReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}

		var id int32
		query := "INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse, text, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
		err = ctx.Conn.QueryRow(context.Background(), query, reference.Book, reference.StartChapter, reference.StartVerse, reference.EndChapter, reference.EndVerse, r.FormValue("text"), *(session.user_id)).Scan(&id)

		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
