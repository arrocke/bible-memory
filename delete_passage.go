package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func DeletePassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		query := "DELETE FROM passage WHERE id = $1"
		_, err = ctx.Conn.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		if strings.HasSuffix(r.Header.Get("Hx-Current-Url"), fmt.Sprintf("/passages/%d/review", id)) {
			w.Header().Set("Hx-Location", "/passages")
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("Delete")
}
