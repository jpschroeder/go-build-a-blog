package main

import (
	"database/sql"
	"html/template"
	"path"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Parse all of the html templates so that they can be rendered with data
func parseTemplates(tmplPath string) (*template.Template, error) {
	return template.ParseGlob(path.Join(tmplPath, "*.html"))
}

// Open and return the sqlite database file
func openDb(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}

// Create the schema in the database if it doesn't already exist
func createSchema(db *sql.DB) error {
	sql := `
		create table if not exists blogs (
			BlogId integer primary key autoincrement,
			Slug varchar(64) not null,
			KeyHash varchar(128) not null,
			IsDefault integer not null,
			Title varchar(64) not null,
			Body text null,
			Html test null
		);

		create unique index if not exists idx_blogs_slug on blogs(Slug);

		create table if not exists pages (
			PageId integer primary key autoincrement,
			BlogId integer REFERENCES blogs(BlogId),
			Slug varchar(64) not null,
			Date datetime not null,
			Show integer not null,
			Title varchar(64) not null,
			Body text null,
			Html text null
		);

		create unique index if not exists idx_pages_slug on pages(BlogId, Slug);
	`
	_, err := db.Exec(sql)
	return err
}

// Open the database and create its schema
func initDb(dbPath string) (*sql.DB, error) {
	db, err := openDb(dbPath)
	if err != nil {
		return nil, err
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Return a router than has all of the handlers registered
func registerRoutes(db *sql.DB, tmpl *template.Template) *mux.Router {
	blogSlug := "/{blogslug:[a-z0-9-]+}"
	pageSlug := "/{pageslug:[a-z0-9-]+}"

	r := mux.NewRouter()
	r.HandleFunc("/", DefaultBlogHandler(db, tmpl)).Methods("GET")

	r.HandleFunc(blogSlug, ViewBlogHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+"/add", AddPageHandler(tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+"/add", CreatePageHandler(db, tmpl)).Methods("POST")

	r.HandleFunc(blogSlug+pageSlug, ViewPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/edit", EditPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/edit", UpdatePageHandler(db, tmpl)).Methods("POST")
	r.HandleFunc(blogSlug+pageSlug+"/delete", DeletePageHandler(tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/delete", DeletePageConfirmHandler(db)).Methods("POST")
	return r
}
