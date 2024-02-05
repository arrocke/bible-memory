package main

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageListItem struct {
	Id        int32
	Reference string
	ReviewAt  string
	ReviewDue bool
}
type PassagesTemplateData struct {
	Passages  []PassageListItem
	StartOpen bool
	LayoutTemplateData
}

func LoadPassagesTemplateData(conn *pgxpool.Pool, user_id int32, tz_offset int) (*PassagesTemplateData, error) {
	type PassageModel struct {
		Id           int32
		Book         string
		StartChapter int32
		StartVerse   int32
		EndChapter   int32
		EndVerse     int32
		ReviewAt     *time.Time
	}

	query := `
		SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
		FROM passage
		WHERE user_id = $1
		ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
	`
	rows, _ := conn.Query(context.Background(), query, user_id)
	defer rows.Close()

	passages, err := pgx.CollectRows(rows, pgx.RowToStructByName[PassageModel])
	if err != nil {
		println(err.Error())
		return nil, err
	}

	layoutTemplateData, err := LoadLayoutTemplateData(conn, &user_id)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	templateData := PassagesTemplateData{
		Passages:           make([]PassageListItem, len(passages)),
		LayoutTemplateData: *layoutTemplateData,
	}
	location := time.FixedZone("temp", tz_offset)
	now := time.Now().In(location)
	for i, passage := range passages {
		passageData := PassageListItem{
			Id:        passage.Id,
			Reference: FormatReference(passage.Book, passage.StartChapter, passage.StartVerse, passage.EndChapter, passage.EndVerse),
		}
		if passage.ReviewAt != nil {
			passageData.ReviewAt = passage.ReviewAt.Format("01-02-2006")
			passageData.ReviewDue = now.Compare(*passage.ReviewAt) > 0
		}
		templateData.Passages[i] = passageData
	}

	return &templateData, nil
}

func GetPassages(router *mux.Router, ctx *ServerContext) {
	tmpl := template.Must(template.ParseFiles("templates/passage_list_partial.html", "templates/passages.html", "templates/layout.html"))

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

		data, err := LoadPassagesTemplateData(ctx.Conn, *session.user_id, GetTZ(r))
		data.StartOpen = true
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout.html", data)
	}).Methods("GET")
}
