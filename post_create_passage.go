package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostCreatePassage(router *mux.Router, conn *pgxpool.Pool) {
	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		reference, err := ParseReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse, text) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err = conn.Exec(context.Background(), query, reference.Book, reference.StartChapter, reference.StartVerse, reference.EndChapter, reference.EndVerse, r.FormValue("Text"))
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods("Post")
}
