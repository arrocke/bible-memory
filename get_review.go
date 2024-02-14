package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func GetPassageReview(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Text         string
		ReviewedAt   *time.Time
		Interval     *int
	}

	type ReviewWord struct {
		Number      string
		Prefix      string
		Word        string
		Suffix      string
		FirstLetter string
		RestOfWord  string
	}
	type PartialTemplateData struct {
		Id              int32
		Reference       string
		Words           []ReviewWord
		AlreadyReviewed bool
		HardInterval    int
		GoodInterval    int
		EasyInterval    int
	}
	type TemplateData struct {
		PartialTemplateData
		PassagesTemplateData
	}

	wordRegex := regexp.MustCompile(`(?:(\d+)\s?)?([^A-Za-zÀ-ÖØ-öø-ÿ\s]+)?([A-Za-zÀ-ÖØ-öø-ÿ]+(?:(?:'|’|-)[A-Za-zÀ-ÖØ-öø-ÿ]+)?(?:'|’)?)([^A-Za-zÀ-ÖØ-öø-ÿ0-9]*\s+)?`)

	parseWords := func(text string) []ReviewWord {
		matches := wordRegex.FindAllStringSubmatch(text, -1)

		words := make([]ReviewWord, len(matches))

		for i, match := range matches {
			words[i] = ReviewWord{
				Number:      match[1],
				Prefix:      match[2],
				Word:        match[3],
				Suffix:      match[4],
				FirstLetter: match[3][0:1],
				RestOfWord:  match[3][1:],
			}
		}

		return words
	}

	tmpl := template.Must(template.ParseFiles("templates/review_partial.html", "templates/review.html", "templates/passage_list_partial.html", "templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages/{Id}/{Mode}", func(w http.ResponseWriter, r *http.Request) {
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

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, reviewed_at, interval FROM passage WHERE id = $1 AND user_id = $2"
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

		now := GetClientDate(r)

		partialTemplateData := PartialTemplateData{
			Id:              passage.Id,
			Reference:       FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
			Words:           parseWords(passage.Text),
			AlreadyReviewed: passage.ReviewedAt != nil && passage.ReviewedAt.Equal(now),
			HardInterval:    GetNextInterval(now, 2, passage.Interval, passage.ReviewedAt),
			GoodInterval:    GetNextInterval(now, 3, passage.Interval, passage.ReviewedAt),
			EasyInterval:    GetNextInterval(now, 4, passage.Interval, passage.ReviewedAt),
		}

		if r.Header.Get("Hx-Current-Url") == "" {
			passagesTemplateData, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetClientDate(r))
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
				partialTemplateData,
				*passagesTemplateData,
			})
		} else {
			tmpl.ExecuteTemplate(w, "review_partial.html", partialTemplateData)
		}
	}).Methods("Get")
}
