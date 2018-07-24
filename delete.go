package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (s Server) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	var i interface{}
	return s.tmpl.ExecuteTemplate(w, "delete.html", i)
}

func (s Server) deleteConfirmHandler(w http.ResponseWriter, r *http.Request) error {
	if !s.verifyKey(r.FormValue("key")) {
		return errors.New("invalid key")
	}

	slug := mux.Vars(r)["slug"]
	err := s.deleteCommand(slug)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func (s Server) deleteCommand(slug string) error {
	sql := `
		delete from pages where Slug = ?
	`
	_, err := s.db.Exec(sql, slug)
	return err
}
