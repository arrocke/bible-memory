package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var wordRegex = regexp.MustCompile(`(\d+ [^A-Za-z0-9']*)?(\w+(?:(?:'|’|-)\w+)?(?:'|’)?)([^A-Za-z0-9']+)?`)

type ReviewWord struct {
	Prefix      string
	Word        string
	Gap         string
	FirstLetter string
	RestOfWord  string
}

func ParseWords(text string) []ReviewWord {
	matches := wordRegex.FindAllStringSubmatch(text, -1)

	words := make([]ReviewWord, len(matches))

	for i, match := range matches {
		words[i] = ReviewWord{
			Prefix:      match[1],
			Word:        match[2],
			Gap:         match[3],
			FirstLetter: match[2][0:1],
			RestOfWord:  match[2][1:],
		}
	}

	return words
}

func GetPassageReview(router *mux.Router, conn *pgxpool.Pool) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Text         string
	}

	type TemplateData struct {
		Id        int32
		Reference string
		Words     []ReviewWord
		Mode      string
	}

	tmpl := template.Must(template.ParseFiles("templates/review.html", "templates/layout.html"))

	router.HandleFunc("/passages/{Id}/{Mode}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text FROM passage WHERE id = $1"
		rows, _ := conn.Query(context.Background(), query, id)
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

		words := ParseWords(passage.Text)

		tmpl.ExecuteTemplate(w, "layout.html", TemplateData{
			Id:        passage.Id,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
			Words:     words,
			Mode:      vars["Mode"],
		})
	}).Methods("Get")
}
