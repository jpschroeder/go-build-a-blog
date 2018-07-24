package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Application")

	//hash, _ := hashPassword("xxx")
	//fmt.Println(hash)

	s, err := initServer()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", s.registerRoutes())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Application Finished")
}
