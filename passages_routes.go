package main

import (
	"context"
	"errors"
	"fmt"
	"main/domain_model"
	"main/view"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
)

func (ctx ServerContext) RenderPassagePage(w http.ResponseWriter, r *http.Request, subpage templ.Component) error {
	userId := GetUserId(r)
	clientDate := GetClientDate(r)

	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

	query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
        FROM passage
        WHERE user_id = $1
        ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
    `
	rows, _ := ctx.Conn.Query(r.Context(), query, userId)
	dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
	if err != nil {
		return err
	}

	passages := make([]view.PassageListItemModel, len(dbpassages))
	for i, dbpassage := range dbpassages {
		passageData := view.PassageListItemModel{
			Id: dbpassage.Id,
			Reference: view.PassageReference{
				Book:         dbpassage.Book,
				StartChapter: dbpassage.StartChapter,
				StartVerse:   dbpassage.StartVerse,
				EndChapter:   dbpassage.EndChapter,
				EndVerse:     dbpassage.EndVerse,
			}.String(),
		}
		if dbpassage.ReviewAt != nil {
			passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
			passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
		}
		passages[i] = passageData
	}

	model := view.PassagesPageModel{
		Passages: passages,
	}

	page := view.PassagesPage(model, subpage)

	return ctx.RenderPage(w, r, page)
}

func parseReference(str string) (domain_model.PassageReference, error) {
	if str == "" {
		return domain_model.PassageReference{},
			errors.New("Please provide a passage reference.")
	}

	parsedReference, err := view.ParsePassageReference(str)
	if err != nil {
		return domain_model.PassageReference{},
			errors.New("Please provide a passage reference in the format Book 1:1-2:5.")
	}
	reference, err := domain_model.NewPassageReference(
		parsedReference.Book,
		parsedReference.StartChapter,
		parsedReference.StartVerse,
		parsedReference.EndChapter,
		parsedReference.EndVerse,
	)
	if err != nil {
		return domain_model.PassageReference{},
			errors.New("Please provide a passage reference in the format Book 1:1-2:5.")
	}

	return reference, nil
}

func (ctx *ServerContext) passagesRoutes(router *http.ServeMux) {
	router.Handle("GET /passages", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		return ctx.RenderPassagePage(w, r, nil)
	})))

	router.Handle("POST /passages", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)

		model := view.AddPassageFormModel{
			Reference: r.FormValue("reference"),
			Text:      r.FormValue("text"),
		}
		hasError := false

        if model.Text == "" {
            model.TextError = "Please provide the passage text."
            hasError = true
        }

        reference, err := parseReference(model.Reference)
        if err != nil {
            model.ReferenceError = err.Error()
            hasError = true
        }

        if hasError {
            return view.AddPassageForm(model).Render(r.Context(), w)
        }

		passage := domain_model.NewPassage(domain_model.NewPassageProps{
			Reference: reference,
			Text:      model.Text,
			Owner:     userId,
		})

		if err := ctx.PassageRepo.Commit(&passage); err != nil {
			return err
		}

        w.Header().Set("Hx-Location", fmt.Sprintf("{\"path\":\"/passages/%d/review\",\"target\":\"#page\"}", passage.Id()))
		w.WriteHeader(http.StatusNoContent)

		return nil
	})))

	router.Handle("GET /passages/new", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		page := view.AddPassagePage(view.AddPassageFormModel{})
		return ctx.RenderPassagePage(w, r, page)
	})))

	router.Handle("GET /passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)

		id, err := ParseInt(r.PathValue("Id"))
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

		type passageModel struct {
			Id           int
			Book         string
			StartChapter int
			StartVerse   int
			EndChapter   int
			EndVerse     int
			Text         string
			ReviewAt     *time.Time
			Interval     *int
		}

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, review_at, interval FROM passage WHERE id = $1 AND user_id = $2"
		rows, _ := ctx.Conn.Query(context.Background(), query, id, userId)
		defer rows.Close()

		passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[passageModel])
		if err != nil {
			return err
		}

		model := view.EditPassagePageModel{
			Id: passage.Id,
			Reference: view.PassageReference{
				Book:         passage.Book,
				StartChapter: passage.StartChapter,
				StartVerse:   passage.StartVerse,
				EndChapter:   passage.EndChapter,
				EndVerse:     passage.EndVerse,
			}.String(),
			Text:     passage.Text,
			Interval: passage.Interval,
			ReviewAt: passage.ReviewAt,
		}

		page := view.EditPassagePage(model)

		return ctx.RenderPassagePage(w, r, page)
	})))

	router.Handle("PUT /passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		id, err := ParseInt(r.PathValue("Id"))
		if err != nil {
			return err
		}

		passage, err := ctx.PassageRepo.Get(id)
		if err != nil {
			return err
		}

		reference, err := parseReference(r.FormValue("reference"))
		if err != nil {
			return err
		}

		var nextReview *domain_model.PassageReview

		interval, err := ParseOptional(ParseInt, r.FormValue("interval"))
		if err != nil {
			return err
		}

		reviewAt, err := ParseOptional(ParseDate, r.FormValue("review_at"))
		if err != nil {
			return err
		}

		if interval != nil && reviewAt != nil {
			nextInterval, err := domain_model.NewReviewInterval(*interval)
			if err != nil {
				return err
			}

			nextReviewAt := domain_model.NewReviewTimestamp(*reviewAt)

			reviewState := passage.Props().ReviewState.Overwrite(nextInterval, nextReviewAt)
			nextReview = &reviewState
		}

		passage.SetReference(reference)
		passage.SetText(r.FormValue("text"))
		passage.SetReviewState(nextReview)

		if err := ctx.PassageRepo.Commit(&passage); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", fmt.Sprintf("/passages/%d/review", id))
		w.WriteHeader(http.StatusNoContent)

		return nil
	})))

	router.Handle("DELETE /passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)

		id, err := ParseInt(r.PathValue("Id"))
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return nil
		}

		query := "DELETE FROM passage WHERE id = $1 AND user_id = $2"
		_, err = ctx.Conn.Exec(context.Background(), query, id, userId)
		if err != nil {
			return err
		}

		if strings.HasSuffix(r.Header.Get("Hx-Current-Url"), fmt.Sprintf("/passages/%d/review", id)) {
			w.Header().Set("Hx-Location", "/passages")
		}
		w.WriteHeader(http.StatusOK)

		return nil
	})))

	router.Handle("GET /passages/{Id}/{Mode}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)
		clientDate := GetClientDate(r)

		id, err := ParseInt(r.PathValue("Id"))
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

		type passageModel struct {
			Id           int
			Book         string
			StartChapter int
			StartVerse   int
			EndChapter   int
			EndVerse     int
			Text         string
			ReviewedAt   *time.Time
			Interval     *int
		}

		query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, reviewed_at, interval FROM passage WHERE id = $1 AND user_id = $2"
		rows, _ := ctx.Conn.Query(r.Context(), query, id, userId)
		defer rows.Close()

		passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[passageModel])
		if err != nil {
			return err
		}

		model := view.ReviewPassagePageModel{
			Id: passage.Id,
			Reference: view.PassageReference{
				Book:         passage.Book,
				StartChapter: passage.StartChapter,
				StartVerse:   passage.StartVerse,
				EndChapter:   passage.EndChapter,
				EndVerse:     passage.EndVerse,
			}.String(),
			Text:            passage.Text,
			AlreadyReviewed: passage.ReviewedAt != nil && passage.ReviewedAt.Equal(clientDate),
			/*
			   HardInterval:    GetNextInterval(now, 2, passage.Interval, passage.ReviewedAt),
			   GoodInterval:    GetNextInterval(now, 3, passage.Interval, passage.ReviewedAt),
			   EasyInterval:    GetNextInterval(now, 4, passage.Interval, passage.ReviewedAt),
			*/
		}

		page := view.ReviewPassagePage(model)

		return ctx.RenderPassagePage(w, r, page)
	})))

	router.Handle("POST /passages/{Id}/review", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		id, err := ParseInt(r.PathValue("Id"))
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

		if r.FormValue("mode") != "review" {
			return view.ReviewResult(view.ReviewResultModel{
				// ReviewAt:
			}).Render(r.Context(), w)
		}

		parsedGrade, err := ParseInt(r.FormValue("grade"))
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return nil
		}

		tz := GetClientTZ(r)

		passage, err := ctx.PassageRepo.Get(id)
		if err != nil {
			return err
		}

		grade, err := domain_model.NewReviewGrade(parsedGrade)
		if err != nil {
			return err
		}

		timestamp := domain_model.NewReviewTimestampForToday(tz)

		passage.Review(grade, timestamp)

		err = ctx.PassageRepo.Commit(&passage)
		if err != nil {
			return err
		}

		return view.ReviewResult(view.ReviewResultModel{
			// ReviewAt:
		}).Render(r.Context(), w)
	})))
}
