package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"time"
)

func AddPageHandler(tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		page := Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
		return tmpl.ExecuteTemplate(w, "edit.html", page)
	})
}

func CreatePageHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		if !verifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}

		page, err1 := parseForm(r)
		if err1 != nil {
			return err1
		}
		slug, err2 := CreatePageCommand(db, page)
		if err2 != nil {
			return err2
		}
		http.Redirect(w, r, "/"+slug+"/edit", http.StatusFound)
		return nil
	})
}

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
