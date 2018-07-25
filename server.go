package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	db   *sql.DB
	tmpl *template.Template
}

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

func initServer(resetHash bool) (Server, error) {
	db, err := openDb()
	if err != nil {
		return Server{}, err
	}
	tmpl, err := parseTemplates()
	if err != nil {
		return Server{}, err
	}
	err = createSchema(db)
	if err != nil {
		return Server{}, err
	}

	s := Server{db: db, tmpl: tmpl}

	if resetHash {
		s.deleteHashCommand()
	}

	err = s.ensureHashExists()
	if err != nil {
		return Server{}, err
	}

	return s, nil
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

type errorHandler func(http.ResponseWriter, *http.Request) error

func makeHandler(fn errorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s Server) registerRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", makeHandler(s.listHandler)).Methods("GET")
	r.HandleFunc("/add", makeHandler(s.addHandler)).Methods("GET")
	r.HandleFunc("/add", makeHandler(s.createHandler)).Methods("POST")
	slugUrl := "/{slug:[a-z0-9-]+}"
	r.HandleFunc(slugUrl, makeHandler(s.viewHandler)).Methods("GET")
	r.HandleFunc(slugUrl+"/edit", makeHandler(s.editHandler)).Methods("GET")
	r.HandleFunc(slugUrl+"/edit", makeHandler(s.updateHandler)).Methods("POST")
	r.HandleFunc(slugUrl+"/delete", makeHandler(s.deleteHandler)).Methods("GET")
	r.HandleFunc(slugUrl+"/delete", makeHandler(s.deleteConfirmHandler)).Methods("POST")
	return r
}
