package main

import (
	"bufio"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Depado/bfchroma"
	"github.com/avelino/slugify"
	"github.com/microcosm-cc/bluemonday"
	stripmd "github.com/writeas/go-strip-markdown"
	"golang.org/x/crypto/bcrypt"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	dateFormat = "2006-01-02"
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
	highlighter := bfchroma.NewRenderer(bfchroma.Style("vs"))
	html := blackfriday.Run(toUnix(input), blackfriday.WithRenderer(highlighter))
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("style").Matching(regexp.MustCompile(`color:#[0-9a-f]+`)).OnElements("span")
	sanitized := policy.SanitizeBytes(html)
	return sanitized
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
	Html  []byte
}

func (p Page) FormattedDate() string {
	return p.Date.Format(dateFormat)
}

// Parse form values into a page object
func parseForm(r *http.Request) (*Page, error) {
	date, err := time.Parse(dateFormat, r.FormValue("date"))
	if err != nil {
		return nil, err
	}

	body := r.FormValue("body")
	//summary := lineNum(body, 1)

	return &Page{
		Date:  date,
		Title: parseTitle(body),
		Body:  []byte(body),
		Show:  r.FormValue("show") == "1"}, nil
}

// Parse the title from the markdown body
func parseTitle(body string) string {
	plain := stripmd.Strip(body)
	title := lineNum(plain, 0)
	return truncate(title, 80)
}

// Get the 0-indexed non-empty line number
func lineNum(in string, num int) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	i := 0
	line := in
	for scanner.Scan() {
		line = scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		if i == num {
			return line
		}
		i++
	}
	return ""
}

// Truncate a string to num of characters
func truncate(str string, num int) string {
	if len(str) <= num {
		return str
	}
	return str[0:num]
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
