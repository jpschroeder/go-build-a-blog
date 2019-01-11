package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the delete page
func DeletePageHandler(tmpl *template.Template) http.HandlerFunc {
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
		return tmpl.ExecuteTemplate(w, "deletepage.html", dto)
	})
}

// Post handler to actually delete a page
func DeletePageConfirmHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]
		key := r.FormValue("key")

		blog, err := BlogMetaQuery(db, blogslug)
		if err != nil {
			return err
		}

		if !verifyHash(key, blog.KeyHash) {
			return errors.New("invalid key")
		}

		err = DeletePageCommand(db, blog.BlogId, pageslug)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/"+blogslug, http.StatusFound)
		return nil
	})
}

// Delete a page from the database
func DeletePageCommand(db *sql.DB, blogId int, pageslug string) error {
	sql := `
		delete from pages where BlogId = ? and Slug = ?
	`
	_, err := db.Exec(sql, blogId, pageslug)
	return err
}
