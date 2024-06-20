package main

import (
	"main/internal/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load(".env")

    app.Start(app.ServerConfiguration {
        DatabaseUrl: os.Getenv("DATABASE_URL"),
        SessionKey: os.Getenv("SESSION_KEY"),
    })
}
