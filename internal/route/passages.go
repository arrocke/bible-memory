package route

import (
	"main/internal/middleware"
	"main/view"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h Handlers) renderPassagesPage(c echo.Context, content templ.Component) error {
    userId := c.Get("user_id").(int)
    clientDate := time.Now() // TODO: get from cookie

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
	rows, _ := h.dbConn.Query(c.Request().Context(), query, userId)
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

	page := view.PassagesPage(model, content)

    app := view.Page(view.AppModel{}, page)
    return app.Render(c.Request().Context(), c.Response().Writer)
}

func (h Handlers) passages(g *echo.Group) {
    g.GET("passages", middleware.AuthMiddleware(true)(func(c echo.Context) error {
        return h.renderPassagesPage(c, nil)
    }))
}
