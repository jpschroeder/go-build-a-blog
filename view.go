package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Date  time.Time
	Show  bool
	Title string
	Body  []byte
}

func (p Page) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

func (p Page) FormattedDateTime() string {
	return p.Date.Format(dateTimeFormat)
}

func ViewPageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		slug := mux.Vars(r)["slug"]

		page, err := ViewPageQuery(db, slug)
		if err != nil {
			return err
		}

		type PageDto struct {
			FormattedDate string
			Title         string
			Body          template.HTML
		}

		body := template.HTML(parseMarkdown(page.Body))
		dto := PageDto{
			FormattedDate: page.FormattedDate(),
			Title:         page.Title,
			Body:          body}

		return tmpl.ExecuteTemplate(w, "view.html", dto)
	})
}

func ViewPageQuery(db *sql.DB, slug string) (*Page, error) {
	sql := `
		select Date, Show, Title, Body from pages where Slug = ?
	`
	row := db.QueryRow(sql, slug)

	var date time.Time
	var show bool
	var title string
	var body []byte
	err := row.Scan(&date, &show, &title, &body)
	if err != nil {
		return nil, err
	}
	return &Page{Date: date, Show: show, Title: title, Body: body}, nil
}
