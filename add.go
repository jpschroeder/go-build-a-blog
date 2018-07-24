package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/avelino/slugify"
)

func (s Server) addHandler(w http.ResponseWriter, r *http.Request) error {
	page := Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
	return s.tmpl.ExecuteTemplate(w, "edit.html", page)
}

func (s Server) createHandler(w http.ResponseWriter, r *http.Request) error {
	if !s.verifyKey(r.FormValue("key")) {
		return errors.New("invalid key")
	}

	page, err1 := parseForm(r)
	if err1 != nil {
		return err1
	}
	slug, err2 := s.createCommand(page)
	if err2 != nil {
		return err2
	}
	http.Redirect(w, r, "/"+slug+"/edit", http.StatusFound)
	return nil
}

func (s Server) createCommand(p *Page) (string, error) {
	sql := `
		insert into pages(Slug, Date, Show, Title, Body) values(?, ?, ?, ?, ?)
	`
	slug := slugify.Slugify(p.Title)
	_, err := s.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body)
	if err != nil {
		return "", err
	}
	return slug, nil
}
