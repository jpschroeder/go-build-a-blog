package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Get http handler that clears the session and redirects the user
func LockBlogHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]

		// Try to read the session cookie
		c, err := r.Cookie(cookieName(blogslug))
		if err != nil {
			http.Redirect(w, r, "/"+blogslug, http.StatusFound)
			return nil
		}

		// Get the token from the cookie
		token := c.Value

		// Remove the session token from the database
		RemoveSessionCommand(db, blogslug, token)

		// Expire the session cookie
		SetCookies(w, clearCookies(blogslug))

		http.Redirect(w, r, "/"+blogslug, http.StatusFound)
		return nil
	})
}

// Get clear session cookies for a root and sub directories
func clearCookies(blogslug string) []*http.Cookie {
	root := clearCookie(blogslug)
	root.Path = fmt.Sprintf("/%s", blogslug)
	sub := clearCookie(blogslug)
	sub.Path = fmt.Sprintf("/%s/", blogslug)
	return []*http.Cookie{root, sub}
}

// The expired cookie value that will clear the cookie from the users browser
func clearCookie(blogslug string) *http.Cookie {
	return &http.Cookie{
		Name:    cookieName(blogslug),
		Value:   "",
		Expires: time.Unix(0, 0),
	}
}

// Database command to remove a session token from the database
func RemoveSessionCommand(db *sql.DB, blogslug string, token string) error {
	sql := `
		delete from sessions 
		where BlogSlug = ? and Token = ?
	`
	_, err := db.Exec(sql, blogslug, token)
	return err
}
