package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type ServerContext struct {
	Conn         *pgxpool.Pool
	SessionStore *sessions.CookieStore
}

type FlashMessage struct {
	Type string
	Text string
}

var FlashTemplate = template.Must(template.ParseFiles("templates/flash.html"))

func main() {
	godotenv.Load(".env")

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	ctx := &ServerContext{
		Conn:         conn,
		SessionStore: store,
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

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	fmt.Printf("Server started on port %d\n", 8080)

	http.ListenAndServe(":8080", r)
}
