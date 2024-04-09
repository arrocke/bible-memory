package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/new", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		session, err := GetSession(r, ctx)
		if err != nil {
			return err
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

        engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
		if r.Header.Get("Hx-Current-Url") == "" {
            if err := engine.RenderCreatePassage(int(*session.user_id), GetClientDate(r)); err != nil {
                return err
            }
		} else {
            if err := engine.RenderCreatePassagePartial(); err != nil {
                return err
            }
		}

        return nil
	})).Methods("Get")
}
