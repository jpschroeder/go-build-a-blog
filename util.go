package main

import (
	"regexp"

	"github.com/avelino/slugify"
	"golang.org/x/crypto/bcrypt"
	blackfriday "gopkg.in/russross/blackfriday.v2"
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

func makeSlug(input string) string {
	return slugify.Slugify(input)
}

func parseMarkdown(input []byte) []byte {
	return blackfriday.Run(toUnix(input))
}

func stripChar(body string, char string) string {
	r := regexp.MustCompile(char)
	return string(r.ReplaceAll([]byte(body), []byte{}))
}

func toUnix(body []byte) []byte {
	r := regexp.MustCompile(`\r`)
	return r.ReplaceAll(body, []byte{})
}
