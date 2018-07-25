package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

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

func UpdatePageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		if !verifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}

		oldSlug := mux.Vars(r)["slug"]

		page, err1 := parseForm(r)
		if err1 != nil {
			return err1
		}
		newSlug, err2 := UpdatePageCommand(db, oldSlug, page)
		if err2 != nil {
			return err2
		}
		if oldSlug != newSlug {
			http.Redirect(w, r, "/"+newSlug+"/edit", http.StatusFound)
			return nil
		} else {
			return tmpl.ExecuteTemplate(w, "edit.html", page)
		}
	})
}

func UpdatePageCommand(db *sql.DB, oldSlug string, p *Page) (string, error) {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?
		where Slug = ?
	`
	slug := makeSlug(p.Title)
	_, err := db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, oldSlug)
	if err != nil {
		return "", err
	}
	return slug, nil
}
