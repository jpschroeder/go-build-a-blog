package main

import (
	"flag"
	"log"
	"net/http"
)

// A *very* simple blogging engine
func main() {
	log.Println("Starting Application")

	// Accept a command line flag -reset
	// This flag allows you to change the key needed to edit/delete posts
	reset := flag.Bool("reset", false, "reset the key used to edit/delete")
	flag.Parse()

	// Connect to the sqlite database and make sure the schema exists
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	// Optionally clear out the old edit key and ask the user for a new one
	if *reset {
		DeleteHashCommand(db)
	}
	EnsureHashExists(db)

	// Parse any html templates used by the application
	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatal(err)
	}

	// Register all of the routing handlers
	http.Handle("/", registerRoutes(db, tmpl))

	// Start the application server
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Application Finished")
}
