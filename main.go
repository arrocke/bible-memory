package main

import (
	"context"
	"fmt"
	"main/db"
	"main/view"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ParseInt(str string) (int, error) {
    val, err := strconv.ParseInt(str, 10, 32)
    if err != nil {
        return 0, err
    }
    return int(val), nil
}

func ParseDate(str string) (time.Time, error) {
    val, err := time.Parse("2006-01-02", str)
    if err != nil {
        return time.Time{}, nil
    }
    return val, nil
}

func ParseOptional[T any](conv func(string) (T, error), str string) (*T, error) {
    if str == "" {
        return nil, nil
    } else {
        val, err := conv(str)
        if err != nil {
            return nil, err
        } else {
            return &val, nil
        }
    }
}

type ServerContext struct {
	Conn           *pgxpool.Pool
	SessionManager *SessionManager
	UserRepo       db.UserRepo
	PassageRepo    db.PassageRepo
}

func (ctx ServerContext) RenderPage(w http.ResponseWriter, r *http.Request, page templ.Component) error {
	if r.Header.Get("hx-target") == "page" {
		return page.Render(r.Context(), w)
	} else {
		session := GetSession(r)

		model := view.AppModel{}

		if session.UserId != nil {
			query := `SELECT first_name, last_name FROM "user" WHERE id = $1`
			rows, _ := ctx.Conn.Query(r.Context(), query, session.UserId)
			user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[view.UserModel])
			if err != nil {
				return err
			}

			model.User = &user
		}

		app := view.Page(model, page)
		return app.Render(r.Context(), w)
	}
}

var decoder = schema.NewDecoder()

func HandleErrors(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			fmt.Printf("Server Error: %v\n", err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}
	})
}

func main() {
	godotenv.Load(".env")

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	passageRepo := db.CreatePgPassageRepo(conn)
	userRepo := db.CreatePgUserRepo(conn)

	ctx := &ServerContext{
		Conn:           conn,
		SessionManager: CreateSessionManager(store),
		UserRepo:       &userRepo,
		PassageRepo:    &passageRepo,
	}

	r := http.NewServeMux()

	ctx.indexRoutes(r)
	ctx.userRoutes(r)
	ctx.passagesRoutes(r)

	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))

	env := os.Getenv("ENV")
	if env == "production" {
		r.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "public, max-age=31536000")
			fs.ServeHTTP(w, r)
		})
	} else {
		r.Handle("GET /assets/", fs)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on port %v\n", port)

	if err := http.ListenAndServe(":"+port, ctx.SessionManager.SessionMiddleware(r)); err != nil {
		fmt.Printf("Error in server: %v", err.Error())
	}
}
