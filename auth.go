package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func verifyKey(db *sql.DB, key string) bool {
	hash, err := GetHashQuery(db)
	if err != nil {
		return false
	}
	return verifyHash(key, hash)
}

func ensureHashExists(db *sql.DB) error {
	if HashExistsQuery(db) {
		return nil
	}
	key, err := promptForKey()
	if err != nil {
		return err
	}
	hash, err := createHash(key)
	if err != nil {
		return err
	}
	return UpdateHashCommand(db, hash)
}

func HashExistsQuery(db *sql.DB) bool {
	hash, err := GetHashQuery(db)
	return err == nil && len(hash) > 0
}

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

func UpdateHashCommand(db *sql.DB, hash string) error {
	sql := `
		insert into config(ConfigId, KeyHash) values(1, ?)
		on conflict(ConfigId) do update set KeyHash=?
	`
	_, err := db.Exec(sql, hash, hash)
	return err
}

func DeleteHashCommand(db *sql.DB) error {
	sql := `
		delete from config
	`
	_, err := db.Exec(sql)
	return err
}

func promptForKey() (string, error) {
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
