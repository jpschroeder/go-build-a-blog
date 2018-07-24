package main

import (
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

func (s Server) listHandler(w http.ResponseWriter, r *http.Request) error {
	pages, err := s.listQuery()
	if err != nil {
		return err
	}
	return s.tmpl.ExecuteTemplate(w, "list.html", pages)
}

func (s Server) listQuery() ([]PageListing, error) {
	var ret []PageListing
	sql := `
		select Slug, Title, Date from pages where Show = 1 order by Date desc
	`
	rows, err := s.db.Query(sql)
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
