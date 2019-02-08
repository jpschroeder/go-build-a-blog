package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Get handler to render the edit page with existing data
func EditPageHandler(db *sql.DB, tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]
		dto, err := EditPageQuery(db, blogslug, pageslug)
		if err != nil {
			return err
		}

		return tmpl(w, "editpage.html", dto)
	})
}

// Post handler to save updated page data
func UpdatePageHandler(db *sql.DB, tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]

		unlocked := IsUnlocked(db, w, r, blogslug)
		if !unlocked {
			http.Redirect(w, r, fmt.Sprintf("/%s/unlock", blogslug), http.StatusFound)
			return nil
		}

		page, err := parseForm(r)
		if err != nil {
			return err
		}

		dto := MapEditPageDto(page, blogslug, pageslug)

		newpageslug, err := UpdatePageCommand(db, blogslug, pageslug, page)
		if err != nil {
			dto.Error = err.Error()
			return tmpl(w, "editpage.html", dto)
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/%s", blogslug, newpageslug), http.StatusFound)
		return nil
	})
}

// Model used to populate the edit page
type EditPageDto struct {
	Title         string
	FormattedDate string
	Show          bool
	Body          []byte
	BlogSlug      string
	PageSlug      string
	Error         string
}

// Generate the Edit page dto
func MapEditPageDto(page *Page, blogslug string, pageslug string) EditPageDto {
	return EditPageDto{
		Title:         page.Title,
		FormattedDate: page.FormattedDate(),
		Show:          page.Show,
		Body:          page.Body,
		BlogSlug:      blogslug,
		PageSlug:      pageslug}
}

// Get the full page data from the database
func EditPageQuery(db *sql.DB, blogslug string, pageslug string) (*EditPageDto, error) {
	sql := `
		select Date, Show, Title, Body 
		from pages
		where BlogSlug = ? and PageSlug = ?
	`
	row := db.QueryRow(sql, blogslug, pageslug)

	var date time.Time
	var show bool
	var title string
	var body []byte
	err := row.Scan(&date, &show, &title, &body)
	if err != nil {
		return nil, err
	}
	return &EditPageDto{
		FormattedDate: date.Format(dateFormat),
		Title:         title,
		Show:          show,
		Body:          body,
		BlogSlug:      blogslug,
		PageSlug:      pageslug,
	}, nil
}

// Update page data in the database
func UpdatePageCommand(db *sql.DB, blogslug string, oldpageslug string, p *Page) (string, error) {
	sql := `
		update pages
		set PageSlug = ?, Date = ?, Show = ?, Title = ?, Body = ?, Html = ?
		where BlogSlug = ? and PageSlug = ?
	`
	newpageslug := makeSlug(p.Title)
	html := parseMarkdown(p.Body)
	_, err := db.Exec(sql, newpageslug, p.Date, p.Show, p.Title, p.Body, html, blogslug, oldpageslug)
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			err = errors.New("There is already a page with this title")
		}
		return "", err
	}
	return newpageslug, nil
}
