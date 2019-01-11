package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Get handler to display the default blog and pages from it
func DefaultBlogHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blog, err := DefaultBlogQuery(db)
		if err != nil {
			return err
		}

		pages, err := ListPagesQuery(db, blog.BlogId)
		if err != nil {
			return err
		}
		model := BlogViewModel{Blog: blog, Pages: pages}
		return tmpl.ExecuteTemplate(w, "viewblog.html", model)
	})
}

// Get handler to display a specific blog and the pages from it
func ViewBlogHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		blog, err := ViewBlogQuery(db, blogslug)
		if err != nil {
			return err
		}

		pages, err := ListPagesQuery(db, blog.BlogId)
		if err != nil {
			return err
		}
		model := BlogViewModel{Blog: blog, Pages: pages}
		return tmpl.ExecuteTemplate(w, "viewblog.html", model)
	})
}

type BlogViewModel struct {
	Blog  *ViewBlogDto
	Pages []PageListing
}

// Listing for a page without full content
type PageListing struct {
	Slug  string
	Title string
	Date  time.Time
}

func (p PageListing) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

// Query the database for the list of page titles and metadata
func ListPagesQuery(db *sql.DB, blogId int) ([]PageListing, error) {
	var ret []PageListing
	sql := `
		select Slug, Title, Date from pages where Show = 1 and BlogId = ? order by Date desc
	`
	rows, err := db.Query(sql, blogId)
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		var title string
		var date time.Time
		err = rows.Scan(&slug, &title, &date)
		if err != nil {
			return ret, err
		}
		ret = append(ret, PageListing{Slug: slug, Title: title, Date: date})
	}
	return ret, rows.Err()
}

// Full blog content
type ViewBlogDto struct {
	BlogId int
	Slug   string
	Title  string
	Html   template.HTML
}

// Get the full blog data from the database
func ViewBlogQuery(db *sql.DB, blogslug string) (*ViewBlogDto, error) {
	sql := `
		select BlogId, Slug, Title, Html from blogs where Slug = ?
	`
	row := db.QueryRow(sql, blogslug)
	return ParseBlogResult(row)
}

// Get the data for the default blog from the database
func DefaultBlogQuery(db *sql.DB) (*ViewBlogDto, error) {
	sql := `
		select BlogId, Slug, Title, Html from blogs where IsDefault = 1
	`
	row := db.QueryRow(sql)
	return ParseBlogResult(row)
}

// Parse a returned sql row into a blog struct
func ParseBlogResult(row *sql.Row) (*ViewBlogDto, error) {
	var blogId int
	var slug string
	var title string
	var html []byte
	err := row.Scan(&blogId, &slug, &title, &html)
	if err != nil {
		return nil, err
	}
	return &ViewBlogDto{
		BlogId: blogId,
		Slug:   slug,
		Title:  title,
		Html:   template.HTML(html)}, nil
}
