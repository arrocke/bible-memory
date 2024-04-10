package main

import (
	"main/services"
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PostReviewPassage(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages/{Id}/review", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

		if r.FormValue("mode") != "review" {
            return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(userId, GetClientDate(r))
		}

		grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid grade", http.StatusBadRequest)
			return nil
		}

		tz := GetClientTZ(r)

		if err := ctx.PassageService.Review(services.ReviewPassageRequest{
			Id:    int(id),
			Grade: int(grade),
			Tz:    tz,
		}); err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return nil
		}

        return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderReviewResult(userId, GetClientDate(r))
	}))).Methods("Post")
}
