package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
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

func (s Server) viewHandler(w http.ResponseWriter, r *http.Request) error {
	slug := mux.Vars(r)["slug"]

	page, err := s.viewQuery(slug)
	if err != nil {
		return err
	}

	type PageDto struct {
		FormattedDate string
		Title         string
		Body          template.HTML
	}

	body := template.HTML(blackfriday.Run(toUnix(page.Body)))
	dto := PageDto{
		FormattedDate: page.FormattedDate(),
		Title:         page.Title,
		Body:          body}

	return s.tmpl.ExecuteTemplate(w, "view.html", dto)
}

func (s Server) viewQuery(slug string) (*Page, error) {
	sql := `
		select Date, Show, Title, Body from pages where Slug = ?
	`
	row := s.db.QueryRow(sql, slug)

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
