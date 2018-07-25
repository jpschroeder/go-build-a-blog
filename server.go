package main

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func parseTemplates() (*template.Template, error) {
	return template.ParseGlob("templates/*.html")
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
		create table if not exists config (
			ConfigId integer primary key autoincrement,
			KeyHash varchar(128) not null
		);
	`
	_, err := db.Exec(sql)
	return err
}

func initDb(resetHash bool) (*sql.DB, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	if resetHash {
		DeleteHashCommand(db)
	}

	err = ensureHashExists(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func parseForm(r *http.Request) (*Page, error) {
	date, err := time.Parse(dateTimeFormat, r.FormValue("date"))
	if err != nil {
		return nil, err
	}
	return &Page{
		Date:  date,
		Title: r.FormValue("title"),
		Body:  []byte(r.FormValue("body")),
		Show:  r.FormValue("show") == "1"}, nil
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func handleErrors(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func requireKey(db *sql.DB, fn handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if !verifyKey(db, r.FormValue("key")) {
			return errors.New("invalid key")
		}
		return fn(w, r)
	}
}

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
