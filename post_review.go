package main

import (
	"main/domain_model"
	"main/services"
	"main/view"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var REVIEW_MAP = [...]int{1, 1, 1, 2, 2, 3, 5, 8, 13, 21, 34, 55}

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

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

		if r.FormValue("mode") != "review" {
			query := `
                SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
                FROM passage
                WHERE user_id = $1
                ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
            `
			rows, _ := ctx.Conn.Query(r.Context(), query, session.user_id)
			dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
			if err != nil {
				http.Error(w, "Error", http.StatusBadRequest)
				return
			}

			clientDate := GetClientDate(r)

			passages := make([]view.PassageListItemModel, len(dbpassages))
			for i, dbpassage := range dbpassages {
				passageData := view.PassageListItemModel{
					Id:        dbpassage.Id,
					Reference: domain_model.PassageReference{dbpassage.Book, dbpassage.StartChapter, dbpassage.StartVerse, dbpassage.EndChapter, dbpassage.EndVerse}.String(),
				}
				if dbpassage.ReviewAt != nil {
					passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
					passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
				}
				passages[i] = passageData
			}

			view.ReviewResult(view.ReviewResultModel{
				Passages: passages,
			}).Render(r.Context(), w)
			return
		}

		grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return
		}

		tz := GetClientTZ(r)

		if err := ctx.PassageService.Review(services.ReviewPassageRequest{
			Id:    int(id),
			Grade: int(grade),
			Tz:    tz,
		}); err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		query := `
            SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
            FROM passage
            WHERE user_id = $1
            ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
        `
		rows, _ := ctx.Conn.Query(r.Context(), query, session.user_id)
		dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
		if err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		clientDate := GetClientDate(r)

		passages := make([]view.PassageListItemModel, len(dbpassages))
		for i, dbpassage := range dbpassages {
			passageData := view.PassageListItemModel{
				Id:        dbpassage.Id,
				Reference: domain_model.PassageReference{dbpassage.Book, dbpassage.StartChapter, dbpassage.StartVerse, dbpassage.EndChapter, dbpassage.EndVerse}.String(),
			}
			if dbpassage.ReviewAt != nil {
				passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
				passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
			}
			passages[i] = passageData
		}

		view.ReviewResult(view.ReviewResultModel{
			// ReviewAt:             passage.ReviewState.NextReview.Value().Format("01-02-2006"),
			Passages: passages,
		}).Render(r.Context(), w)
	}).Methods("Post")
}
