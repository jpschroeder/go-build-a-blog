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

	s, err := initServer(*reset)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", s.registerRoutes())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Application Finished")
}
