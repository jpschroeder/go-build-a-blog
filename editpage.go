package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Get handler to render the edit page with existing data
func EditPageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]
		page, err := ViewPageQuery(db, blogslug, pageslug)
		if err != nil {
			return err
		}

		dto := MapEditPageDto(page, blogslug, pageslug)
		return tmpl.ExecuteTemplate(w, "editpage.html", dto)
	})
}

// Post handler to save updated page data
func UpdatePageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]
		key := r.FormValue("key")

		blog, err := BlogMetaQuery(db, blogslug)
		if err != nil {
			return err
		}

		page, err := parseForm(r)
		if err != nil {
			return err
		}

		dto := MapEditPageDto(page, blogslug, pageslug)

		if !verifyHash(key, blog.KeyHash) {
			dto.Error = "invalid key"
			return tmpl.ExecuteTemplate(w, "editpage.html", dto)
		}

		newSlug, err := UpdatePageCommand(db, blog.BlogId, pageslug, page)
		if err != nil {
			dto.Error = err.Error()
			return tmpl.ExecuteTemplate(w, "editpage.html", dto)
		}

		if pageslug != newSlug {
			http.Redirect(w, r, fmt.Sprintf("/%s/%s/edit", blogslug, newSlug), http.StatusFound)
			return nil
		} else {
			return tmpl.ExecuteTemplate(w, "editpage.html", dto)
		}
	})
}

// Model used to populate the edit page
type EditPageDto struct {
	Title             string
	FormattedDateTime string
	Show              bool
	Body              []byte
	BlogSlug          string
	PageSlug          string
	Error             string
}

// Generate the Edit page dto
func MapEditPageDto(page *Page, blogslug string, pageslug string) EditPageDto {
	return EditPageDto{
		Title:             page.Title,
		FormattedDateTime: page.FormattedDateTime(),
		Show:              page.Show,
		Body:              page.Body,
		BlogSlug:          blogslug,
		PageSlug:          pageslug}
}

// Update page data in the database
func UpdatePageCommand(db *sql.DB, blogId int, oldSlug string, p *Page) (string, error) {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?, Html = ?
		where BlogId = ? and Slug = ?
	`
	slug := makeSlug(p.Title)
	html := parseMarkdown(p.Body)
	_, err := db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, html, blogId, oldSlug)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: pages.BlogId, pages.Slug" {
			err = errors.New("There is already a page with this title")
		}
		return "", err
	}
	return slug, nil
}
