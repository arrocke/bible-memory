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
	router.Handle("/passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return nil
		}

		query := "DELETE FROM passage WHERE id = $1 AND user_id = $2"
		_, err = ctx.Conn.Exec(context.Background(), query, id, userId)
		if err != nil {
			return err
		}

		if strings.HasSuffix(r.Header.Get("Hx-Current-Url"), fmt.Sprintf("/passages/%d/review", id)) {
			w.Header().Set("Hx-Location", "/passages")
		}
		w.WriteHeader(http.StatusOK)

        return nil
	}))).Methods("Delete")
}
