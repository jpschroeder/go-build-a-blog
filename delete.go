package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func DeletePageHandler(tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		var i interface{}
		return tmpl.ExecuteTemplate(w, "delete.html", i)
	})
}

func DeletePageConfirmHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		if !verifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}

		slug := mux.Vars(r)["slug"]
		err := DeletePageCommand(db, slug)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	})
}

func DeletePageCommand(db *sql.DB, slug string) error {
	sql := `
		delete from pages where Slug = ?
	`
	_, err := db.Exec(sql, slug)
	return err
}
