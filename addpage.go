package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Get handler to render the edit page with empty data
func AddPageHandler(tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		page := Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
		return tmpl.ExecuteTemplate(w, "editpage.html", page)
	})
}

// Post handler to save a new page
func CreatePageHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		key := r.FormValue("key")

		blog, err := BlogMetaQuery(db, blogslug)
		if err != nil {
			return err
		}

		if !verifyHash(key, blog.KeyHash) {
			return errors.New("invalid key")
		}

		page, err := parseForm(r)
		if err != nil {
			return err
		}

		pageslug, err := CreatePageCommand(db, blog.BlogId, page)
		if err != nil {
			return err
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/%s/edit", blogslug, pageslug), http.StatusFound)
		return nil
	})
}

// metadata about a blog
type BlogMeta struct {
	BlogId  int
	KeyHash string
}

// Get the metadata for a specified blog
func BlogMetaQuery(db *sql.DB, blogslug string) (*BlogMeta, error) {
	sql := `
		select BlogId, KeyHash from blogs where Slug = ?
	`
	row := db.QueryRow(sql, blogslug)

	var blogId int
	var hash string
	err := row.Scan(&blogId, &hash)
	if err != nil {
		return nil, err
	}
	return &BlogMeta{BlogId: blogId, KeyHash: hash}, nil
}

// Insert a new page into the database
func CreatePageCommand(db *sql.DB, blogId int, p *Page) (string, error) {
	sql := `
		insert into pages(BlogId, Slug, Date, Show, Title, Body, Html) values(?, ?, ?, ?, ?, ?, ?)
	`
	pageslug := makeSlug(p.Title)
	html := parseMarkdown(p.Body)
	_, err := db.Exec(sql, blogId, pageslug, p.Date, p.Show, p.Title, p.Body, html)
	if err != nil {
		return "", err
	}
	return pageslug, nil
}
