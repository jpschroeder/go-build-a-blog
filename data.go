// Data Repository

package main

import (
	"database/sql"
	"time"

	"github.com/avelino/slugify"
	_ "github.com/mattn/go-sqlite3"
)

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02T15:04"

type PageListing struct {
	Slug  string
	Title string
	Date  time.Time
}

func (p PageListing) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

type Page struct {
	Date  time.Time
	Show  bool
	Title string
	Body  []byte
}

func (p Page) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

func (p Page) FormattedDateTime() string {
	return p.Date.Format(dateTimeFormat)
}

type Data struct {
	db *sql.DB
}

func openDb() (*sql.DB, error) {
	return sql.Open("sqlite3", "data.db")
}

func createSchema(db *sql.DB) error {
	sql := `
		create table if not exists pages (
			PageId integer primary key autoincrement,
			Slug varchar(64) not null,
			Date datetime not null,
			Show integer not null,
			Title varchar(64) not null,
			Body text null
		);
		create unique index if not exists idx_pages_slug on pages(Slug);
	`
	_, err := db.Exec(sql)
	return err
}

func (r *Data) init() error {
	db, err := openDb()
	if err != nil {
		return err
	}
	r.db = db
	return createSchema(r.db)
}

func (r Data) list() ([]PageListing, error) {
	var ret []PageListing
	sql := `
		select Slug, Title, Date from pages where Show = 1 order by Date desc
	`
	rows, err := r.db.Query(sql)
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

func (r Data) create(p *Page) (string, error) {
	sql := `
		insert into pages(Slug, Date, Show, Title, Body) values(?, ?, ?, ?, ?)
	`
	slug := slugify.Slugify(p.Title)
	_, err := r.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (r Data) update(oldSlug string, p *Page) (string, error) {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?
		where Slug = ?
	`
	slug := slugify.Slugify(p.Title)
	_, err := r.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, oldSlug)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (r Data) view(slug string) (*Page, error) {
	sql := `
		select Date, Show, Title, Body from pages where Slug = ?
	`
	row := r.db.QueryRow(sql, slug)

	var date time.Time
	var show bool
	var title string
	var body []byte
	err := row.Scan(&date, &show, &title, &body)
	if err != nil {
		return nil, err
	}
	return &Page{Date: date, Show: show, Title: title, Body: body}, nil
}

func (r Data) delete(slug string) error {
	sql := `
		delete from pages where Slug = ?
	`
	_, err := r.db.Exec(sql, slug)
	return err
}
