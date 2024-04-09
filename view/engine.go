package view

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ViewEngine struct {
    conn *pgxpool.Pool
    context context.Context
    writer http.ResponseWriter
}

func CreateViewEngine(conn *pgxpool.Pool, context context.Context, w http.ResponseWriter) ViewEngine {
    return ViewEngine {
        conn: conn,
        context: context,
        writer: w,
    }
}

func (eng ViewEngine) LoadAppModel(user_id *int, page interface{}) (AppModel,error) {
    model := AppModel {
        Page: page,
    }

    if user_id != nil {
        query := `SELECT first_name, last_name FROM "user" WHERE id = $1`
        rows, _ := eng.conn.Query(eng.context, query, user_id)
        user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[UserModel])
        if err != nil {
            return model,err
        }

        model.User = &user
    }

    return model, nil
}

func (eng ViewEngine) LoadPassagesPageModel(user_id int, clientDate time.Time, page interface{}) (PassagesPageModel, error) {
    model := PassagesPageModel {
        Page: page,
    }

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
    rows, _ := eng.conn.Query(context.Background(), query, user_id)
    dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
    if err != nil {
        return model,err
    }

    model.Passages = make([]PassageListItemModel, len(dbpassages))
    for i, dbpassage := range dbpassages {
        passageData := PassageListItemModel{
            Id:        dbpassage.Id,
            Reference: PassageReference{
                Book: dbpassage.Book, 
                StartChapter: dbpassage.StartChapter,
                StartVerse: dbpassage.StartVerse,
                EndChapter: dbpassage.EndChapter,
                EndVerse: dbpassage.EndVerse,
            }.String(),
        }
        if dbpassage.ReviewAt != nil {
            passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
            passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
        }
        model.Passages[i] = passageData
    }

    return model, nil
}

func (eng ViewEngine) LoadProfilePageModel(user_id int) (ProfilePageModel, error) {
    query := `SELECT email, first_name, last_name FROM "user" WHERE id = $1`
    rows, _ := eng.conn.Query(eng.context, query, user_id)
    model, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ProfilePageModel])
    if err != nil {
        return ProfilePageModel{}, err
    }
    
    return model, nil
}

func (eng ViewEngine) LoadEditPassageModel(userId int, passageId int) (EditPassagePageModel, error) {
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
    rows, _ := eng.conn.Query(context.Background(), query, passageId, userId)
    defer rows.Close()

    passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[passageModel])
    if err != nil {
        return EditPassagePageModel{}, nil
    }

    model := EditPassagePageModel {
        Id:        passage.Id,
        Reference: PassageReference{
            Book: passage.Book, 
            StartChapter: passage.StartChapter,
            StartVerse: passage.StartVerse,
            EndChapter: passage.EndChapter,
            EndVerse: passage.EndVerse,
        }.String(),
        Text:      passage.Text,
        Interval:  passage.Interval,
        ReviewAt:  passage.ReviewAt,
    }

    return model, nil
}

func (eng ViewEngine) LoadReviewPassageModel(userId int, passageId int, clientDate time.Time) (ReviewPassagePageModel, error) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		Text         string
		ReviewedAt     *time.Time
		Interval     *int
	}

    query := "SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, reviewed_at, interval FROM passage WHERE id = $1 AND user_id = $2"
    rows, _ := eng.conn.Query(context.Background(), query, passageId, userId)
    defer rows.Close()

    passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[passageModel])
    if err != nil {
        return ReviewPassagePageModel{}, nil
    }

    model := ReviewPassagePageModel{
        Id:              passage.Id,
        Reference: PassageReference{
            Book: passage.Book, 
            StartChapter: passage.StartChapter,
            StartVerse: passage.StartVerse,
            EndChapter: passage.EndChapter,
            EndVerse: passage.EndVerse,
        }.String(),
        Text: passage.Text,
        AlreadyReviewed: passage.ReviewedAt != nil && passage.ReviewedAt.Equal(clientDate),
        /*
        HardInterval:    GetNextInterval(now, 2, passage.Interval, passage.ReviewedAt),
        GoodInterval:    GetNextInterval(now, 3, passage.Interval, passage.ReviewedAt),
        EasyInterval:    GetNextInterval(now, 4, passage.Interval, passage.ReviewedAt),
        */
    }

    return model, nil
}

func (eng ViewEngine) RenderIndex() error {
    return App(AppModel {
        Page: PublicIndexModel{},
    }).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderLogin() error {
    return App(AppModel {
        Page: LoginPageModel{},
    }).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderLoginError(error string, email string) error {
    return App(AppModel {
        Page: LoginPageModel{
            Error: error,
            Email: email,
        },
    }).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderRegister() error {
    return App(AppModel {
        Page: RegisterPageModel{},
    }).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderProfile(user_id int) error {
    page, err := eng.LoadProfilePageModel(user_id)
    if err != nil {
        return err
    }
    
    return App(AppModel {
        Page: page,
        User: &UserModel{
            FirstName: page.FirstName,
            LastName: page.LastName,
        },
    }).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderPassages(user_id int, clientDate time.Time) error {
    page, err := eng.LoadPassagesPageModel(user_id, clientDate, nil)
    if err != nil {
        return err
    }

    model,err := eng.LoadAppModel(&user_id, page)
    if err != nil {
        return err
    }

    return App(model).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderCreatePassagePartial() error {
    return AddPassagePage(AddPassagePageModel{}).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderCreatePassage(userId int, clientDate time.Time) error {
    page, err := eng.LoadPassagesPageModel(userId, clientDate, AddPassagePageModel{})
    if err != nil {
        return err
    }

    model,err := eng.LoadAppModel(&userId, page)
    if err != nil {
        return err
    }

    return App(model).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderPassageEditPartial(userId int, passageId int) error {
    page, err := eng.LoadEditPassageModel(userId, passageId)
    if err != nil {
        return err
    }

    return EditPassagePage(page).Render(eng.context, eng.writer)
}
func (eng ViewEngine) RenderPassageEdit(userId int, passageId int, clientDate time.Time) error {
    subpage, err := eng.LoadEditPassageModel(userId, passageId)
    if err != nil {
        return err
    }

    page, err := eng.LoadPassagesPageModel(userId, clientDate, subpage)
    if err != nil {
        return err
    }

    model,err := eng.LoadAppModel(&userId, page)
    if err != nil {
        return err
    }

    return App(model).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderReviewPassagePartial(userId int, passageId int, clientDate time.Time) error {
    page, err := eng.LoadReviewPassageModel(userId, passageId, clientDate)
    if err != nil {
        return err
    }

    return ReviewPassagePage(page).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderReviewPassage(userId int, passageId int, clientDate time.Time) error {
    subpage, err := eng.LoadReviewPassageModel(userId, passageId, clientDate)
    if err != nil {
        return err
    }

    page, err := eng.LoadPassagesPageModel(userId, clientDate, subpage)
    if err != nil {
        return err
    }

    model,err := eng.LoadAppModel(&userId, page)
    if err != nil {
        return err
    }

    return App(model).Render(eng.context, eng.writer)
}

func (eng ViewEngine) RenderReviewResult(userId int, clientDate time.Time) error {
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
    rows, _ := eng.conn.Query(eng.context, query, userId)
    dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
    if err != nil {
        return err
    }

    passages := make([]PassageListItemModel, len(dbpassages))
    for i, dbpassage := range dbpassages {
        passageData := PassageListItemModel{
            Id:        dbpassage.Id,
            Reference: PassageReference{
                Book: dbpassage.Book, 
                StartChapter: dbpassage.StartChapter,
                StartVerse: dbpassage.StartVerse,
                EndChapter: dbpassage.EndChapter,
                EndVerse: dbpassage.EndVerse,
            }.String(),
        }
        if dbpassage.ReviewAt != nil {
            passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
            passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
        }
        passages[i] = passageData
    }

    return ReviewResult(ReviewResultModel{
        Passages: passages,
        // ReviewAt: 
    }).Render(eng.context, eng.writer)
}
