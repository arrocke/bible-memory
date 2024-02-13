package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PutEditPassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		reference, err := ParseReference(r.FormValue("reference"))
		if err != nil {
			http.Error(w, "Invalid reference", http.StatusBadRequest)
			return
		}

		reviewAtString := r.FormValue("review_at")
		var reviewAt *string
		if r.FormValue("review_at") != "" {
			reviewAt = &reviewAtString
		}

		query := "UPDATE passage SET book = $3, start_chapter = $4, start_verse = $5, end_chapter = $6, end_verse = $7, text = $8, review_at = $9 WHERE id = $1 AND user_id = $2"
		_, err = ctx.Conn.Exec(context.Background(), query, id, *session.user_id, reference.Book, reference.StartChapter, reference.StartVerse, reference.EndChapter, reference.EndVerse, r.FormValue("text"), reviewAt)
		if err != nil {
			println(err.Error())
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Put")
}
