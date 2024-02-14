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

var REVIEW_MAP = [...]int{1, 1, 1, 2, 2, 3, 5, 8, 13, 21, 34, 55}

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id         int32
		ReviewedAt *time.Time
		ReviewAt   *time.Time
		Interval   *int
	}

	type TemplateData struct {
		Grade    int
		ReviewAt string
		PassagesTemplateData
	}

	tmpl := template.Must(template.ParseFiles("templates/passage_list_partial.html", "templates/review_result.html"))

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

		query := "SELECT id, reviewed_at, review_at, interval FROM passage WHERE id = $1 AND user_id = $2"
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

		_grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 32)
		if err != nil || _grade < 0 || _grade > 4 {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return
		}
		grade := int(_grade)

		if r.FormValue("mode") != "review" {
			passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{PassagesTemplateData: *passagesTemplateData})
			return
		}

		now := GetClientDate(r)

		if passage.ReviewedAt != nil && passage.ReviewedAt.Equal(now) {
			passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{ReviewAt: passage.ReviewAt.Format("01-02-2006"), PassagesTemplateData: *passagesTemplateData})
			return
		}

		interval := GetNextInterval(now, grade, passage.Interval, passage.ReviewedAt)
		newDate := now.AddDate(0, 0, interval)

		query = "UPDATE passage SET review_at = $2, reviewed_at = $3, interval = $4 WHERE id = $1"
		_, err = ctx.Conn.Exec(context.Background(), query, id, newDate, now, interval)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{Grade: grade, ReviewAt: newDate.Format("01-02-2006"), PassagesTemplateData: *passagesTemplateData})
	}).Methods("Post")
}
