package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the delete page
func DeletePageHandler(tmpl ExecuteTemplateFunc) http.HandlerFunc {
	type DeletePageDto struct {
		BlogSlug string
		PageSlug string
	}
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]
		dto := DeletePageDto{
			BlogSlug: blogslug,
			PageSlug: pageslug}
		return tmpl(w, "deletepage.html", dto)
	})
}

// Post handler to actually delete a page
func DeletePageConfirmHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]

		unlocked := IsUnlocked(db, w, r, blogslug)
		if !unlocked {
			http.Redirect(w, r, fmt.Sprintf("/%s/unlock", blogslug), http.StatusFound)
			return nil
		}

		err := DeletePageCommand(db, blogslug, pageslug)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/"+blogslug, http.StatusFound)
		return nil
	})
}

// Delete a page from the database
func DeletePageCommand(db *sql.DB, blogslug string, pageslug string) error {
	sql := `
		delete from pages where BlogSlug = ? and PageSlug = ?
	`
	_, err := db.Exec(sql, blogslug, pageslug)
	return err
}
