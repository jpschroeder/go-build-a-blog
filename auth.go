package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// Check a key against the hash stored in the database
func VerifyKey(db *sql.DB, key string) bool {
	hash, err := GetHashQuery(db)
	if err != nil {
		return false
	}
	return verifyHash(key, hash)
}

// Check if a security hash exists in the database
func HashExistsQuery(db *sql.DB) bool {
	hash, err := GetHashQuery(db)
	return err == nil && len(hash) > 0
}

// Query for the security hash in the database
func GetHashQuery(db *sql.DB) (string, error) {
	sql := `
		select KeyHash from config order by ConfigId desc
	`
	row := db.QueryRow(sql)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, err
}

// Add or update the security hash in the database
func UpdateHashCommand(db *sql.DB, hash string) error {
	sql := `
		insert into config(ConfigId, KeyHash) values(1, ?)
		on conflict(ConfigId) do update set KeyHash=excluded.KeyHash
	`
	_, err := db.Exec(sql, hash)
	return err
}

// Clear the security hash from the database
func DeleteHashCommand(db *sql.DB) error {
	sql := `
		delete from config
	`
	_, err := db.Exec(sql)
	return err
}

// If a hash does not exist in the database
// Prompt the user for one on the command line and store it
func EnsureHashExists(db *sql.DB) error {
	if HashExistsQuery(db) {
		return nil
	}
	key, err := PromptForKey()
	if err != nil {
		return err
	}
	hash, err := createHash(key)
	if err != nil {
		return err
	}
	return UpdateHashCommand(db, hash)
}

// Prompt the user for a key on the command line and return it
func PromptForKey() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter key: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	key = stripChar(stripChar(key, `\n`), `\r`)
	if len(key) < 1 {
		return "", errors.New("invalid key")
	}
	return key, nil
}
