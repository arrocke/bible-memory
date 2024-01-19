package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := mux.NewRouter()

	GetPassages(r, conn)
	GetPassageReview(r, conn)
	GetCreatePassage(r, conn)
	PostCreatePassage(r, conn)
	DeletePassage(r, conn)

	http.ListenAndServe(":8080", r)
}
