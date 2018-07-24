package main

import (
	"errors"
	"net/http"

	"github.com/avelino/slugify"
	"github.com/gorilla/mux"
)

func (s Server) editHandler(w http.ResponseWriter, r *http.Request) error {
	slug := mux.Vars(r)["slug"]
	page, err := s.viewQuery(slug)
	if err != nil {
		return err
	}
	return s.tmpl.ExecuteTemplate(w, "edit.html", page)
}

func (s Server) updateHandler(w http.ResponseWriter, r *http.Request) error {
	if !s.verifyKey(r.FormValue("key")) {
		return errors.New("invalid key")
	}

	oldSlug := mux.Vars(r)["slug"]

	page, err1 := parseForm(r)
	if err1 != nil {
		return err1
	}
	newSlug, err2 := s.updateCommand(oldSlug, page)
	if err2 != nil {
		return err2
	}
	if oldSlug != newSlug {
		http.Redirect(w, r, "/"+newSlug+"/edit", http.StatusFound)
		return nil
	} else {
		return s.tmpl.ExecuteTemplate(w, "edit.html", page)
	}
}

func (s Server) updateCommand(oldSlug string, p *Page) (string, error) {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?
		where Slug = ?
	`
	slug := slugify.Slugify(p.Title)
	_, err := s.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, oldSlug)
	if err != nil {
		return "", err
	}
	return slug, nil
}
