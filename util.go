package main

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02T15:04"

func createHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func verifyHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func stripChar(body string, char string) string {
	r := regexp.MustCompile(char)
	return string(r.ReplaceAll([]byte(body), []byte{}))
}

func toUnix(body []byte) []byte {
	r := regexp.MustCompile(`\r`)
	return r.ReplaceAll(body, []byte{})
}
