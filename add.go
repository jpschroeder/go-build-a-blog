package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"time"
)

// Get handler to render the edit page with empty data
func AddPageHandler(tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		page := Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
		return tmpl.ExecuteTemplate(w, "edit.html", page)
	})
}

// Post handler to save a new page
func CreatePageHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		if !VerifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}

		page, err := parseForm(r)
		if err != nil {
			return err
		}

		slug, err := CreatePageCommand(db, page)
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/"+slug+"/edit", http.StatusFound)
		return nil
	})
}

// Insert a new page into the database
func CreatePageCommand(db *sql.DB, p *Page) (string, error) {
	sql := `
		insert into pages(Slug, Date, Show, Title, Body) values(?, ?, ?, ?, ?)
	`
	slug := makeSlug(p.Title)
	_, err := db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body)
	if err != nil {
		return "", err
	}
	return slug, nil
}
