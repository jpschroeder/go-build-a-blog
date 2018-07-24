package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func (s Server) verifyKey(key string) bool {
	hash, err := s.hashQuery()
	if err != nil {
		return false
	}
	return verifyHash(key, hash)
}

func (s Server) ensureHashExists() error {
	if s.hashExists() {
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
	return s.updateHashCommand(hash)
}

func (s Server) hashExists() bool {
	hash, err := s.hashQuery()
	return err != nil && len(hash) > 0
}

func (s Server) hashQuery() (string, error) {
	sql := `
		select KeyHash from config order by ConfigId desc
	`
	row := s.db.QueryRow(sql)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (s Server) updateHashCommand(hash string) error {
	sql := `
		insert into config(ConfigId, KeyHash) values(1, ?)
		on conflict(ConfigId) do update set KeyHash=?
	`
	_, err := s.db.Exec(sql, hash, hash)
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
