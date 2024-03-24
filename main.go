package main

import (
	"context"
	"fmt"
	"html/template"
	"main/db"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
    "github.com/gorilla/schema"
)

type ServerContext struct {
	Conn         *pgxpool.Pool
	SessionStore *sessions.CookieStore
    PassageRepo  db.PassageRepo
}

type FlashMessage struct {
	Type string
	Text string
}

var FlashTemplate = template.Must(template.ParseFiles("templates/flash.html"))

var decoder = schema.NewDecoder()

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

	ctx := &ServerContext{
		Conn:         conn,
		SessionStore: store,
        PassageRepo: &passageRepo,
	}

	r := mux.NewRouter()

	GetIndex(r, ctx)

	GetRegister(r, ctx)
	GetLogin(r, ctx)
	GetProfile(r, ctx)
	PostRegister(r, ctx)
	PostLogin(r, ctx)
	PostLogout(r, ctx)
	PutProfile(r, ctx)

	GetPassages(r, ctx)
	GetPassageReview(r, ctx)
	GetCreatePassage(r, ctx)
	GetPassageEdit(r, ctx)
	PostCreatePassage(r, ctx)
	PostReviewPassage(r, ctx)
	PutEditPassage(r, ctx)
	DeletePassage(r, ctx)

	env := os.Getenv("ENV")
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/assets/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env == "production" {
			w.Header().Add("Cache-Control", "public, max-age=31536000")
		}
		fs.ServeHTTP(w, r)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on port %v\n", port)

	http.ListenAndServe(":"+port, r)
}
