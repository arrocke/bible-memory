package main

import (
	"context"
	"fmt"
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

		var id int32
		query := "INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse, text) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
		err = conn.QueryRow(context.Background(), query, reference.Book, reference.StartChapter, reference.StartVerse, reference.EndChapter, reference.EndVerse, r.FormValue("text")).Scan(&id)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
