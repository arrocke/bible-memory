package main

import (
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPassageEdit(router *mux.Router, ctx *ServerContext) {
	router.Handle("/passages/{Id}", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
		if r.Header.Get("Hx-Current-Url") == "" {
            return engine.RenderPassageEdit(userId, int(id), GetClientDate(r))
		} else {
            return engine.RenderPassageEditPartial(userId, int(id))
		}
	}))).Methods("Get")
}
