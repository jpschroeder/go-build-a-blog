package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func stripChar(body string, char string) string {
	r := regexp.MustCompile(char)
	return string(r.ReplaceAll([]byte(body), []byte{}))
}

type Auth struct {
	hash string
}

func (a *Auth) init() {
	authFile := "hash.db"
	storedhash, err := ioutil.ReadFile(authFile)
	if err == nil {
		// auth file exists
		a.hash = string(storedhash)
	} else {
		// auth file doesn't exist
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter key: ")
		key, _ := reader.ReadString('\n')
		key = stripChar(stripChar(key, `\n`), `\r`)
		enteredhash, _ := hashPassword(key)
		ioutil.WriteFile(authFile, []byte(enteredhash), 0644)
		a.hash = enteredhash
	}
}

func (a Auth) checkPassword(password string) bool {
	return checkPasswordHash(password, a.hash)
}
