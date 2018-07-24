package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	db   *sql.DB
	tmpl *template.Template
	hash string
}

func (s Server) checkPassword(password string) bool {
	return checkPasswordHash(password, s.hash)
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
	`
	_, err := db.Exec(sql)
	return err
}

func readHash() (string, error) {
	authFile := "hash.db"
	storedhash, err1 := ioutil.ReadFile(authFile)
	if err1 == nil {
		// auth file exists
		return string(storedhash), nil
	} else {
		// auth file doesn't exist
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter key: ")
		key, err2 := reader.ReadString('\n')
		if err2 != nil {
			return "", err2
		}
		key = stripChar(stripChar(key, `\n`), `\r`)
		enteredhash, err3 := hashPassword(key)
		if err3 != nil {
			return "", err3
		}
		ioutil.WriteFile(authFile, []byte(enteredhash), 0644)
		return enteredhash, nil
	}
}

func initServer() (Server, error) {
	db, err1 := openDb()
	if err1 != nil {
		return Server{}, err1
	}
	err2 := createSchema(db)
	if err2 != nil {
		return Server{}, err2
	}
	hash, err3 := readHash()
	if err3 != nil {
		return Server{}, err3
	}
	tmpl, err4 := parseTemplates()
	if err4 != nil {
		return Server{}, err4
	}
	return Server{
		db:   db,
		tmpl: tmpl,
		hash: hash}, nil
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
