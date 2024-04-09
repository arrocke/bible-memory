package main

import (
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPassageEdit(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
		if r.Header.Get("Hx-Current-Url") == "" {
            if err := engine.RenderPassageEdit(int(*session.user_id), int(id), GetClientDate(r)); err != nil {
                return err
            }
		} else {
            if err := engine.RenderPassageEditPartial(int(*session.user_id), int(id)); err != nil {
                return err
            }
		}

        return nil
	})).Methods("Get")
}
