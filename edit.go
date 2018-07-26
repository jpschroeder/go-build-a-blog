package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the edit page with existing data
func EditPageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		slug := mux.Vars(r)["slug"]
		page, err := ViewPageQuery(db, slug)
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, "edit.html", page)
	})
}

// Post handler to save updated page data
func UpdatePageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		if !VerifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}

		oldSlug := mux.Vars(r)["slug"]

		page, err := parseForm(r)
		if err != nil {
			return err
		}
		newSlug, err := UpdatePageCommand(db, oldSlug, page)
		if err != nil {
			return err
		}
		if oldSlug != newSlug {
			http.Redirect(w, r, "/"+newSlug+"/edit", http.StatusFound)
			return nil
		} else {
			return tmpl.ExecuteTemplate(w, "edit.html", page)
		}
	})
}

// Update page data in the database
func UpdatePageCommand(db *sql.DB, oldSlug string, p *Page) (string, error) {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?, Html = ?
		where Slug = ?
	`
	slug := makeSlug(p.Title)
	html := parseMarkdown(p.Body)
	_, err := db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, html, oldSlug)
	if err != nil {
		return "", err
	}
	return slug, nil
}
