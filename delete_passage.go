package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DeletePassage(router *mux.Router, conn *pgxpool.Pool) {
	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		query := "DELETE FROM passage WHERE id = $1"
		_, err = conn.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("Delete")
}
