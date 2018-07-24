package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Handlers struct {
	data Data
	tmpl *template.Template
}

func (h *Handlers) init() {
	h.tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func (h Handlers) listHandler(w http.ResponseWriter, r *http.Request) error {
	pages, err := h.data.list()
	if err != nil {
		return err
	}
	return h.tmpl.ExecuteTemplate(w, "list.html", pages)
}

func (h Handlers) viewHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	slug := vars["slug"]
	page, err := h.data.view(slug)
	if err != nil {
		return err
	}

	type PageDto struct {
		FormattedDate string
		Title         string
		Body          template.HTML
	}

	toUnix := func(body []byte) []byte {
		r := regexp.MustCompile(`\r`)
		return r.ReplaceAll(body, []byte{})
	}

	body := template.HTML(blackfriday.Run(toUnix(page.Body)))
	dto := PageDto{
		FormattedDate: page.FormattedDate(),
		Title:         page.Title,
		Body:          body}

	return h.tmpl.ExecuteTemplate(w, "view.html", dto)
}

func (h Handlers) addHandler(w http.ResponseWriter, r *http.Request) error {
	page := Page{Date: time.Now(), Title: "", Body: make([]byte, 0), Show: true}
	return h.tmpl.ExecuteTemplate(w, "edit.html", page)
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

func (h Handlers) createHandler(w http.ResponseWriter, r *http.Request) error {
	page, err1 := parseForm(r)
	if err1 != nil {
		return err1
	}
	slug, err2 := h.data.create(page)
	if err2 != nil {
		return err2
	}
	http.Redirect(w, r, "/"+slug+"/edit", http.StatusFound)
	return nil
}

func (h Handlers) editHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	slug := vars["slug"]
	page, err := h.data.view(slug)
	if err != nil {
		return err
	}
	return h.tmpl.ExecuteTemplate(w, "edit.html", page)
}

func (h Handlers) updateHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	oldSlug := vars["slug"]
	page, err1 := parseForm(r)
	if err1 != nil {
		return err1
	}
	newSlug, err2 := h.data.update(oldSlug, page)
	if err2 != nil {
		return err2
	}
	if oldSlug != newSlug {
		http.Redirect(w, r, "/"+newSlug+"/edit", http.StatusFound)
		return nil
	} else {
		return h.tmpl.ExecuteTemplate(w, "edit.html", page)
	}
}

func (h Handlers) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	slug := vars["slug"]
	err := h.data.delete(slug)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
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

func (h Handlers) registerRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", makeHandler(h.listHandler)).Methods("GET")
	r.HandleFunc("/add", makeHandler(h.addHandler)).Methods("GET")
	r.HandleFunc("/add", makeHandler(h.createHandler)).Methods("POST")
	r.HandleFunc("/{slug:[a-z0-9-]+}", makeHandler(h.viewHandler)).Methods("GET")
	r.HandleFunc("/{slug:[a-z0-9-]+}/edit", makeHandler(h.editHandler)).Methods("GET")
	r.HandleFunc("/{slug:[a-z0-9-]+}/edit", makeHandler(h.updateHandler)).Methods("POST")
	r.HandleFunc("/{slug:[a-z0-9-]+}/delete", makeHandler(h.deleteHandler)).Methods("GET")
	return r
}
