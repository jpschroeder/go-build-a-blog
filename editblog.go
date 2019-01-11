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
		blog, err := ViewBlogQuery(db, blogslug)
		if err != nil {
			return err
		}
		dto := EditBlogDto{
			Title:    blog.Title,
			Body:     blog.Body,
			BlogSlug: blog.Slug}

		return tmpl.ExecuteTemplate(w, "editblog.html", dto)
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
