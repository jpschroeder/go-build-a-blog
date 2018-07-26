package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"
)

type PageListing struct {
	Slug  string
	Title string
	Date  time.Time
}

func (p PageListing) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

// Get handler to display a list of pages to be used as an index
func ListPagesHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		pages, err := ListPagesQuery(db)
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, "list.html", pages)
	})
}

// Query the database for the list of page titles and metadata
func ListPagesQuery(db *sql.DB) ([]PageListing, error) {
	var ret []PageListing
	sql := `
		select Slug, Title, Date from pages where Show = 1 order by Date desc
	`
	rows, err := db.Query(sql)
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
