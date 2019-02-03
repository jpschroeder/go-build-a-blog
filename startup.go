package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// If a blog does not exist in the database
// Prompt the user for one on the command line and store it
func EnsureDefaultBlogExists(db *sql.DB) error {
	if DefaultBlogExistsQuery(db) {
		return nil
	}

	fmt.Println("Enter blog title: ")
	title, err := readInput()
	if err != nil {
		return err
	}

	fmt.Println("Enter key: ")
	key, err := readInput()
	if err != nil {
		return err
	}

	return AddDefaultBlogCommand(db, key, title)
}

// Prompt the user for a new default key and update it in the database
func ResetDefaultBlogKey(db *sql.DB) error {
	fmt.Println("Enter new key: ")
	key, err := readInput()
	if err != nil {
		return err
	}

	hash, err := createHash(key)
	if err != nil {
		return err
	}

	return ChangeDefaultBlogKeyCommand(db, hash)
}

// Check if a default blog exists in the database
func DefaultBlogExistsQuery(db *sql.DB) bool {
	sql := `
		select exists(select BlogSlug from blogs where IsDefault = 1)
	`
	row := db.QueryRow(sql)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// Add or a default blog to the database
func AddDefaultBlogCommand(db *sql.DB, key string, title string) error {
	blogslug := makeSlug(title)
	hash, err := createHash(key)
	if err != nil {
		return err
	}

	sql := `
		insert into blogs(BlogSlug, KeyHash, IsDefault, Title, Body, Html) 
		values(?, ?, 1, ?, '', '')
	`
	_, err = db.Exec(sql, blogslug, hash, title)
	return err
}

// Add or a default blog to the database
func ChangeDefaultBlogKeyCommand(db *sql.DB, hash string) error {
	sql := `
		update blogs set KeyHash = ? where IsDefault = 1
	`
	_, err := db.Exec(sql, hash)
	return err
}

// Prompt the user for a string on the command line and return it
func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = stripChar(stripChar(input, `\n`), `\r`)
	if len(input) < 1 {
		return "", errors.New("invalid input")
	}
	return input, nil
}
