package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Get handler that clears the session and redirects the user
func LockBlogHandler(db *sql.DB) http.HandlerFunc {
	return handleErrors(func(w http.ResponseWriter, r *http.Request) error {
		blogslug := mux.Vars(r)["blogslug"]

		c, err := r.Cookie(cookieName(blogslug))
		if err != nil {
			http.Redirect(w, r, "/"+blogslug, http.StatusFound)
			return nil
		}

		token := c.Value

		http.SetCookie(w, clearCookie(blogslug))
		RemoveSessionCommand(db, blogslug, token)

		http.Redirect(w, r, "/"+blogslug, http.StatusFound)
		return nil
	})
}

func RemoveSessionCommand(db *sql.DB, blogslug string, token string) error {
	sql := `
		delete from sessions 
		where BlogSlug = ? and Token = ?
	`
	_, err := db.Exec(sql, blogslug, token)
	return err
}

func clearCookie(blogslug string) *http.Cookie {
	return &http.Cookie{
		Name:    cookieName(blogslug),
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
}
