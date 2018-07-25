package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Application")

	reset := flag.Bool("reset", false, "reset the key")
	flag.Parse()

	db, err := initDb(*reset)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", registerRoutes(db, tmpl))

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Application Finished")
}
