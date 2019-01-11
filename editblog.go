package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the edit blog page
func EditBlogHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		blog, err := EditBlogQuery(db, blogslug)
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, "editblog.html", blog)
	})
}

// Post handler to save updated blog data
func UpdateBlogHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		key := r.FormValue("key")
		body := []byte(r.FormValue("body"))
		title := r.FormValue("title")

		dto := EditBlogDto{
			Title:    title,
			Body:     body,
			BlogSlug: blogslug}

		blog, err := BlogMetaQuery(db, blogslug)
		if err != nil {
			return err
		}

		if !verifyHash(key, blog.KeyHash) {
			dto.Error = "invalid key"
			return tmpl.ExecuteTemplate(w, "editblog.html", dto)
		}

		err = UpdateBlogCommand(db, blog.BlogId, body)
		if err != nil {
			return nil
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
	Error    string
}

// Update blog data in the database
func UpdateBlogCommand(db *sql.DB, blogId int, body []byte) error {
	sql := `
		update blogs
		set Body = ?, Html = ?
		where BlogId = ?
	`
	html := parseMarkdown(body)
	_, err := db.Exec(sql, body, html, blogId)
	return err
}

// Get the full blog data from the database
func EditBlogQuery(db *sql.DB, blogslug string) (*EditBlogDto, error) {
	sql := `
		select BlogId, Slug, Title, Body from blogs where Slug = ?
	`
	row := db.QueryRow(sql, blogslug)
	var blogId int
	var slug string
	var title string
	var body []byte
	err := row.Scan(&blogId, &slug, &title, &body)
	if err != nil {
		return nil, err
	}
	return &EditBlogDto{
		BlogSlug: slug,
		Title:    title,
		Body:     body}, nil
}
