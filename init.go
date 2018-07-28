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
		create table if not exists pages (
			PageId integer primary key autoincrement,
			Slug varchar(64) not null,
			Date datetime not null,
			Show integer not null,
			Title varchar(64) not null,
			Body text null,
			Html text null
		);
		create unique index if not exists idx_pages_slug on pages(Slug);
		create table if not exists config (
			ConfigId integer primary key autoincrement,
			KeyHash varchar(128) not null
		);
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
	r := mux.NewRouter()
	r.HandleFunc("/", ListPagesHandler(db, tmpl)).Methods("GET")
	r.HandleFunc("/add", AddPageHandler(tmpl)).Methods("GET")
	r.HandleFunc("/add", CreatePageHandler(db)).Methods("POST")
	slugUrl := "/{slug:[a-z0-9-]+}"
	r.HandleFunc(slugUrl, ViewPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(slugUrl+"/edit", EditPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(slugUrl+"/edit", UpdatePageHandler(db, tmpl)).Methods("POST")
	r.HandleFunc(slugUrl+"/delete", DeletePageHandler(tmpl)).Methods("GET")
	r.HandleFunc(slugUrl+"/delete", DeletePageConfirmHandler(db)).Methods("POST")
	return r
}
