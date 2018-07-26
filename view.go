package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// A get handler to view the html of a page
func ViewPageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	type PageDto struct {
		FormattedDate string
		Title         string
		Html          template.HTML
	}
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		slug := mux.Vars(r)["slug"]

		page, err := ViewPageQuery(db, slug)
		if err != nil {
			return err
		}

		dto := PageDto{
			FormattedDate: page.FormattedDate(),
			Title:         page.Title,
			Html:          template.HTML(page.Html)}

		return tmpl.ExecuteTemplate(w, "view.html", dto)
	})
}

// Get the full page data from the database
func ViewPageQuery(db *sql.DB, slug string) (*Page, error) {
	sql := `
		select Date, Show, Title, Body, Html from pages where Slug = ?
	`
	row := db.QueryRow(sql, slug)

	var date time.Time
	var show bool
	var title string
	var body []byte
	var html []byte
	err := row.Scan(&date, &show, &title, &body, &html)
	if err != nil {
		return nil, err
	}
	return &Page{Date: date, Show: show, Title: title, Body: body, Html: html}, nil
}
