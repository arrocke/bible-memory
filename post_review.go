package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var REVIEW_MAP = [...]int32{1, 1, 1, 2, 2, 3, 5, 8, 13, 21, 34, 55}

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id    int32
		Level int32
	}

	type TemplateData struct {
		Kind  string
		Level int32
	}

	tmpl := template.Must(template.ParseFiles("templates/review_result.html"))

	router.HandleFunc("/passages/{Id}/review", func(w http.ResponseWriter, r *http.Request) {
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

		query := "SELECT id, level FROM passage WHERE id = $1 AND user_id = $2"
		rows, _ := ctx.Conn.Query(context.Background(), query, id, *session.user_id)
		defer rows.Close()

		passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[PassageModel])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		if r.FormValue("mode") != "review" {
			tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{Kind: "non-review", Level: passage.Level})
			return
		}

		accuracy, err := strconv.ParseFloat(r.FormValue("accuracy"), 64)
		if err != nil {
			http.Error(w, "Invalid accuracy", http.StatusBadRequest)
			return
		}

		var kind string
		var newLevel int32
		if accuracy == 1.0 {
			kind = "perfect"
			newLevel = min(passage.Level+1, int32(len(REVIEW_MAP)))
		} else if accuracy > 0.9 {
			kind = "good"
			newLevel = passage.Level
		} else if accuracy > 0.8 {
			kind = "ok"
			newLevel = max(passage.Level-1, 0)
		} else {
			kind = "fail"
			newLevel = passage.Level / 2
		}

		newDate := time.Now().Add(time.Duration(REVIEW_MAP[newLevel]))

		query = "UPDATE passage SET level = $2, review_at = $3 WHERE id = $1"
		_, err = ctx.Conn.Exec(context.Background(), query, id, newLevel, newDate)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{Kind: kind, Level: passage.Level})

	}).Methods("Post")
}
