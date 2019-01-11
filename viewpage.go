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
	type ViewPageDto struct {
		FormattedDate string
		Title         string
		Html          template.HTML
		BlogSlug      string
		PageSlug      string
	}
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]

		page, err := ViewPageQuery(db, blogslug, pageslug)
		if err != nil {
			return err
		}

		dto := ViewPageDto{
			FormattedDate: page.FormattedDate(),
			Title:         page.Title,
			Html:          template.HTML(page.Html),
			BlogSlug:      blogslug,
			PageSlug:      pageslug}

		return tmpl.ExecuteTemplate(w, "viewpage.html", dto)
	})
}

// Get the full page data from the database
func ViewPageQuery(db *sql.DB, blogslug string, pageslug string) (*Page, error) {
	sql := `
		select p.Date, p.Show, p.Title, p.Body, p.Html 
		from pages as p
		inner join blogs as b on p.BlogId = b.BlogId
		where b.Slug = ? and p.Slug = ?
	`
	row := db.QueryRow(sql, blogslug, pageslug)

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
