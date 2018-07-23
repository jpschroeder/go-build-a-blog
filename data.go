// Data Repository

package main

import (
	"database/sql"
	"time"

	"github.com/avelino/slugify"
	_ "github.com/mattn/go-sqlite3"
)

type PageListing struct {
	Slug  string
	Title string
	Date  time.Time
}

type Page struct {
	Date  time.Time
	Show  bool
	Title string
	Body  string
}

type Data struct {
	db *sql.DB
}

func openDb() *sql.DB {
	filename := "data.db"
	db, err := sql.Open("sqlite3", filename)
	checkErr(err)
	return db
}

func createSchema(db *sql.DB) {
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
	checkErr(err)
}

func (r *Data) init() {
	r.db = openDb()
	createSchema(r.db)
}

func (r Data) list() []PageListing {
	sql := `
		select Slug, Title, Date from pages order by Date desc
	`
	rows, err := r.db.Query(sql)
	checkErr(err)
	defer rows.Close()

	var ret []PageListing
	for rows.Next() {
		var slug string
		var title string
		var date time.Time
		err = rows.Scan(&slug, &title, &date)
		checkErr(err)
		ret = append(ret, PageListing{Slug: slug, Title: title, Date: date})
	}
	checkErr(rows.Err())
	return ret
}

func (r Data) create(p *Page) string {
	sql := `
		insert into pages(Slug, Date, Show, Title, Body) values(?, ?, ?, ?, ?)
	`
	slug := slugify.Slugify(p.Title)
	_, err := r.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body)
	checkErr(err)
	return slug
}

func (r Data) update(oldSlug string, p *Page) string {
	sql := `
		update pages
		set Slug = ?, Date = ?, Show = ?, Title = ?, Body = ?
		where Slug = ?
	`
	slug := slugify.Slugify(p.Title)
	_, err := r.db.Exec(sql, slug, p.Date, p.Show, p.Title, p.Body, oldSlug)
	checkErr(err)
	return slug
}

func (r Data) view(slug string) *Page {
	sql := `
		select Date, Show, Title, Body from pages where Slug = ?
	`
	row := r.db.QueryRow(sql, slug)

	var date time.Time
	var show bool
	var title string
	var body string
	err := row.Scan(&date, &show, &title, &body)
	checkErr(err)
	return &Page{Date: date, Show: show, Title: title, Body: body}
}

func (r Data) delete(slug string) {
	sql := `
		delete from pages where Slug = ?
	`
	_, err := r.db.Exec(sql, slug)
	checkErr(err)
}
