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

func GetPassageReview(router *mux.Router, ctx *ServerContext) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		Text         string
	}

	type PartialTemplateData struct {
		Id        int32
		Reference string
		Words     []ReviewWord
		Mode      string
	}

	type FullTemplateData struct {
		Id        int32
		Reference string
		Words     []ReviewWord
		Mode      string
		Passages  []PassageListItem
	}

	tmpl := template.Must(template.ParseFiles("templates/review_partial.html", "templates/review.html", "templates/passages.html", "templates/layout.html"))

	router.HandleFunc("/passages/{Id}/{Mode}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mode := vars["Mode"]
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text FROM passage WHERE id = $1"
		rows, _ := ctx.Conn.Query(context.Background(), query, id)
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

		if r.Header.Get("Hx-Current-Url") == "" {
			passageList, err := LoadPassageList(ctx.Conn)
			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			tmpl.ExecuteTemplate(w, "layout.html", FullTemplateData{
				Id:        passage.Id,
				Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
				Words:     words,
				Mode:      mode,
				Passages:  passageList,
			})
		} else {
			tmpl.ExecuteTemplate(w, "review_partial.html", PartialTemplateData{
				Id:        passage.Id,
				Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
				Words:     words,
				Mode:      mode,
			})
		}

	}).Methods("Get")
}
