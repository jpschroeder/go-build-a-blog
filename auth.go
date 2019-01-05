package main

import (
	"database/sql"
)

// Check a key against the hash stored in the database
func VerifyKey(db *sql.DB, key string) bool {
	hash, err := defaultBlogKeyQuery(db)
	if err != nil {
		return false
	}
	return verifyHash(key, hash)
}

// Check if a default blog exists in the database
func defaultBlogKeyQuery(db *sql.DB) (string, error) {
	sql := `
		select KeyHash from blogs where IsDefault = 1
	`
	row := db.QueryRow(sql)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}
