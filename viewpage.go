package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ViewPageDto struct {
	FormattedDate string
	Title         string
	Html          template.HTML
	BlogSlug      string
	PageSlug      string
	Unlocked      bool
}

// A get handler to view the html of a page
func ViewPageHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		pageslug := mux.Vars(r)["pageslug"]

		dto, err := ViewPageQuery(db, blogslug, pageslug)
		if err != nil {
			return err
		}

		dto.Unlocked = IsUnlocked(db, w, r, blogslug)

		return tmpl.ExecuteTemplate(w, "viewpage.html", dto)
	})
}

// Get the full page data from the database
func ViewPageQuery(db *sql.DB, blogslug string, pageslug string) (*ViewPageDto, error) {
	sql := `
		select Date, Show, Title, Html 
		from pages
		where BlogSlug = ? and PageSlug = ?
	`
	row := db.QueryRow(sql, blogslug, pageslug)

	var date time.Time
	var show bool
	var title string
	var html []byte
	err := row.Scan(&date, &show, &title, &html)
	if err != nil {
		return nil, err
	}
	return &ViewPageDto{
		FormattedDate: date.Format(dateFormat),
		Title:         title,
		Html:          template.HTML(html),
		BlogSlug:      blogslug,
		PageSlug:      pageslug,
	}, nil
}
