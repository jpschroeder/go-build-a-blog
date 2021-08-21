package main

import (
	"database/sql"
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type ExecuteTemplateFunc func(wr io.Writer, name string, data interface{}) error

//go:embed templates
var templates embed.FS

// Parse all of the html templates from the executable bundle so that they can be rendered with data
func parseBundledTemplates() (ExecuteTemplateFunc, error) {
	tmpl, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		return nil, err
	}
	return tmpl.ExecuteTemplate, nil
}

// Parse all of the html templates from the file system so that they can be rendered with data
// This is used for development mode
func parseFileTemplates() (ExecuteTemplateFunc, error) {
	// Re-parse the templates on every call
	// This allows for "hot-reloading" of the page with template updates
	// It will also be much slower and should not be used for production
	tmplFunc := func(wr io.Writer, name string, data interface{}) error {
		// You could theoretically only parse the particular template being called
		// However, that would not pick up changes to templates that it depends on
		tmpl, err := template.ParseGlob("./templates/*")
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(wr, name, data)
	}
	return tmplFunc, nil
}

// Open and return the sqlite database file
func openDb(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}

// Create the schema in the database if it doesn't already exist
func createSchema(db *sql.DB) error {
	sql := `
		create table if not exists blogs (
			BlogSlug varchar(64) not null primary key,
			KeyHash varchar(128) not null,
			IsDefault integer not null,
			Title varchar(64) not null,
			Body text null,
			Html test null
		);

		create index if not exists idx_blogs_default on blogs(IsDefault) where IsDefault = 1;

		create table if not exists pages (
			PageSlug varchar(64) not null,
			BlogSlug varchar(64) not null references blogs(BlogSlug) on update cascade on delete cascade,
			Date datetime not null,
			Show integer not null,
			Title varchar(64) not null,
			Body text null,
			Html text null,
			primary key (PageSlug, BlogSlug)
		);

		create index if not exists idx_pages_list on pages(BlogSlug, Show, Date);

		create table if not exists sessions (
			Token varchar(64) not null,
			BlogSlug varchar(64) not null references blogs(BlogSlug) on update cascade on delete cascade,
			Effective datetime not null default(datetime('now')),
			primary key (Token, BlogSlug)
		);

		create table if not exists cache (
			Name varchar(128) not null primary key,
			Data text null
		);
	`
	_, err := db.Exec(sql)

	sql = `
		alter table pages add column Summary varchar(512) not null default "";
	`
	// ignore any errors that come out of this script (columns already exist)
	db.Exec(sql)

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

//go:embed static
var static embed.FS

// Return a router than has all of the handlers registered
func registerRoutes(db *sql.DB, tmpl ExecuteTemplateFunc, devmode bool) *mux.Router {
	blogSlug := "/{blogslug:[a-z0-9-]+}"
	pageSlug := "/{pageslug:[a-z0-9-]+}"

	r := mux.NewRouter()
	r.HandleFunc("/", DefaultBlogHandler(db, tmpl)).Methods("GET")

	r.HandleFunc(blogSlug, ViewBlogHandler(db, tmpl)).Methods("GET")

	r.HandleFunc(blogSlug+"/unlock", UnlockBlogHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+"/unlock", DoUnlockBlogHandler(db, tmpl)).Methods("POST")
	r.HandleFunc(blogSlug+"/lock", LockBlogHandler(db)).Methods("GET")

	r.HandleFunc(blogSlug+"/edit", EditBlogHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+"/edit", UpdateBlogHandler(db, tmpl)).Methods("POST")

	r.HandleFunc(blogSlug+"/add", AddPageHandler(tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+"/add", CreatePageHandler(db, tmpl)).Methods("POST")

	r.HandleFunc(blogSlug+pageSlug, ViewPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/edit", EditPageHandler(db, tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/edit", UpdatePageHandler(db, tmpl)).Methods("POST")
	r.HandleFunc(blogSlug+pageSlug+"/delete", DeletePageHandler(tmpl)).Methods("GET")
	r.HandleFunc(blogSlug+pageSlug+"/delete", DeletePageConfirmHandler(db)).Methods("POST")

	if devmode {
		// Read files directly from the /static directory when in dev mode
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	} else {
		// Read files bundled into the executable when in production mode
		r.PathPrefix("/static/").Handler(http.FileServer(http.FS(static)))
	}

	return r
}
