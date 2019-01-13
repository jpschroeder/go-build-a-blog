package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const sessionTimeoutDays = 14

// view model used to render the unlock page
type UnlockBlogDto struct {
	BlogSlug string
	Error    string
}

// Get handler to render the unlock page
func UnlockBlogHandler(tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		dto := UnlockBlogDto{BlogSlug: blogslug}
		return tmpl.ExecuteTemplate(w, "unlock.html", dto)
	})
}

// Post handler to actually unlock a blog
func DoUnlockBlogHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]
		key := r.FormValue("key")
		dto := UnlockBlogDto{BlogSlug: blogslug}

		keyhash, err := BlogKeyQuery(db, blogslug)
		if err != nil {
			return err
		}

		if !verifyHash(key, keyhash) {
			dto.Error = "invalid key"
			return tmpl.ExecuteTemplate(w, "unlock.html", dto)
		}

		token := uuid.NewV4().String()
		err = InsertSessionCommand(db, blogslug, token)
		if err != nil {
			return err
		}

		http.SetCookie(w, sessionCookie(blogslug, token))

		http.Redirect(w, r, "/"+blogslug, http.StatusFound)
		return nil
	})
}

// Get the key hash for a specified blog
func BlogKeyQuery(db *sql.DB, blogslug string) (string, error) {
	sql := `
		select KeyHash from blogs where BlogSlug = ?
	`
	row := db.QueryRow(sql, blogslug)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Insert a new page into the database
func InsertSessionCommand(db *sql.DB, blogslug string, token string) error {
	sql := `
		insert into sessions(BlogSlug, Token) values(?, ?)
	`
	_, err := db.Exec(sql, blogslug, token)
	return err
}

func LookupSessionQuery(db *sql.DB, blogslug string, token string) bool {
	sql := fmt.Sprintf(`
		select exists(
			select token from sessions
			where BlogSlug = ? and Token = ? 
			and Effective > datetime('now', '-%d days')
		)
	`, sessionTimeoutDays)
	row := db.QueryRow(sql, blogslug, token)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func RefreshSessionCommand(db *sql.DB, blogslug string, token string) error {
	sql := `
		update sessions set Effective = datetime('now') 
		where BlogSlug = ? and Token = ? 
	`
	_, err := db.Exec(sql, blogslug, token)
	return err
}

func ExpireSessionsJob(db *sql.DB) {
	for {
		log.Println("Running Expire Sessions Job")
		err := ExpireSessionsCommand(db)
		if err != nil {
			log.Println("Error in Expire Sessions Job")
			log.Println(err)
		}
		time.Sleep(24 * time.Hour)
	}
}

func ExpireSessionsCommand(db *sql.DB) error {
	sql := fmt.Sprintf(`
		delete from sessions 
		where Effective <= datetime('now', '-%d days')
	`, sessionTimeoutDays)
	_, err := db.Exec(sql)
	return err
}

func IsUnlocked(db *sql.DB, w http.ResponseWriter, r *http.Request, blogslug string) bool {
	c, err := r.Cookie(cookieName(blogslug))
	if err != nil {
		return false
	}

	token := c.Value

	exists := LookupSessionQuery(db, blogslug, token)
	if !exists {
		return false
	}

	RefreshSessionCommand(db, blogslug, token)

	http.SetCookie(w, sessionCookie(blogslug, token))
	return true
}

func sessionCookie(blogslug string, token string) *http.Cookie {
	return &http.Cookie{
		Name:    cookieName(blogslug),
		Value:   token,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, sessionTimeoutDays),
	}
}

// Get the name of the session cookie from the blog slug
func cookieName(blogslug string) string {
	return fmt.Sprintf("session_token_%s", blogslug)
}
