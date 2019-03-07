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

// Get handler to render the edit page with empty data
func AddPageHandler(tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		page := &Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
		dto := MapEditPageDto(page, blogslug, "")
		return tmpl(w, "editpage.html", dto)
	})
}

// Post handler to save a new page
func CreatePageHandler(db *sql.DB, tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]

		unlocked := IsUnlocked(db, w, r, blogslug)
		if !unlocked {
			http.Redirect(w, r, fmt.Sprintf("/%s/unlock", blogslug), http.StatusFound)
			return nil
		}

		page, err := parseForm(r)
		if err != nil {
			return err
		}

		dto := MapEditPageDto(page, blogslug, "")

		pageslug, err := CreatePageCommand(db, blogslug, page)
		if err != nil {
			dto.Error = err.Error()
			return tmpl(w, "editpage.html", dto)
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/%s", blogslug, pageslug), http.StatusFound)
		return nil
	})
}

// Insert a new page into the database
func CreatePageCommand(db *sql.DB, blogslug string, p *Page) (string, error) {
	sql := `
		insert into pages(BlogSlug, PageSlug, Date, Show, Title, Body, Html, Summary) 
		values(?, ?, ?, ?, ?, ?, ?, ?)
	`
	pageslug := makeSlug(p.Title)
	html := parseMarkdown(p.Body)
	summary := lineNum(string(p.Body), 1)
	summaryHtml := string(parseMarkdown([]byte(summary)))
	_, err := db.Exec(sql, blogslug, pageslug, p.Date, p.Show, p.Title, p.Body, html, summaryHtml)
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			err = errors.New("There is already a page with this title")
		}
		return "", err
	}
	return pageslug, nil
}
