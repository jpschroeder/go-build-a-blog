package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/avelino/slugify"
	"golang.org/x/crypto/bcrypt"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = "2006-01-02T15:04"
)

// Generate a cryptographic password hash
func createHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// Verify a plaintext password against a hash
func verifyHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Translate a string into a url friendly value
func makeSlug(input string) string {
	return slugify.Slugify(input)
}

// Translate a markdown string into html
func parseMarkdown(input []byte) []byte {
	return blackfriday.Run(toUnix(input))
}

// Remove all instances of a character from a string
func stripChar(body string, char string) string {
	r := regexp.MustCompile(char)
	return string(r.ReplaceAll([]byte(body), []byte{}))
}

// Translate a string with windows style newlines into one with unix style newlines
// i.e. \r\n to just \n
func toUnix(body []byte) []byte {
	r := regexp.MustCompile(`\r`)
	return r.ReplaceAll(body, []byte{})
}

// Represent a page stored in the database
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

// Parse form values into a page object
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

// An http handler function that returns an error
type handlerFunc func(http.ResponseWriter, *http.Request) error

// Translate an http handler function that returns an error a regular http handler
func handleErrors(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
