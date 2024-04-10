package main

import (
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPassageReview(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages/{Id}/{Mode}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
		if r.Header.Get("Hx-Current-Url") == "" {
            return engine.RenderReviewPassage(userId, int(id), GetClientDate(r))
		} else {
            return engine.RenderReviewPassagePartial(userId, int(id), GetClientDate(r))
		}
	}))).Methods("Get")
}
