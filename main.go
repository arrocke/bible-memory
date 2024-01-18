package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	GetPassages(r)
	GetPassageReview(r)

	http.ListenAndServe(":8080", r)
}
