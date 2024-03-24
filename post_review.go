package main

import (
	"html/template"
	"main/domain_model"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

		passage, err := ctx.PassageRepo.Get(uint(id))
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

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

		grade, err := domain_model.ParseReviewGrade(r.FormValue("grade"))
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return
		}

		passage.Review(grade, now)

		err = ctx.PassageRepo.Commit(&passage)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "review_result.html", TemplateData{
			Grade:                int(grade),
			ReviewAt:             passage.ReviewState.NextReview.Value().Format("01-02-2006"),
			PassagesTemplateData: *passagesTemplateData,
		})
	}).Methods("Post")
}
