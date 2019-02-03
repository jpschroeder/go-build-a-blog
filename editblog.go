package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the edit blog page
func EditBlogHandler(db *sql.DB, tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		blog, err := EditBlogQuery(db, blogslug)
		if err != nil {
			return err
		}
		return tmpl(w, "editblog.html", blog)
	})
}

// Post handler to save updated blog data
func UpdateBlogHandler(db *sql.DB, tmpl ExecuteTemplateFunc) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		body := []byte(r.FormValue("body"))

		unlocked := IsUnlocked(db, w, r, blogslug)
		if !unlocked {
			http.Redirect(w, r, fmt.Sprintf("/%s/unlock", blogslug), http.StatusFound)
			return nil
		}

		err := UpdateBlogCommand(db, blogslug, body)
		if err != nil {
			return err
		}

		http.Redirect(w, r, fmt.Sprintf("/%s", blogslug), http.StatusFound)
		return nil
	})
}

// Model used to populate the edit blog page
type EditBlogDto struct {
	Title    string
	Body     []byte
	BlogSlug string
}

// Update blog data in the database
func UpdateBlogCommand(db *sql.DB, blogslug string, body []byte) error {
	sql := `
		update blogs
		set Body = ?, Html = ?
		where BlogSlug = ?
	`
	html := parseMarkdown(body)
	_, err := db.Exec(sql, body, html, blogslug)
	return err
}

// Get the full blog data from the database
func EditBlogQuery(db *sql.DB, blogslug string) (*EditBlogDto, error) {
	sql := `
		select Title, Body from blogs where BlogSlug = ?
	`
	row := db.QueryRow(sql, blogslug)
	var title string
	var body []byte
	err := row.Scan(&title, &body)
	if err != nil {
		return nil, err
	}
	return &EditBlogDto{
		BlogSlug: blogslug,
		Title:    title,
		Body:     body}, nil
}
