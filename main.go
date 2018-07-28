package main

import (
	"flag"
	"log"
	"net/http"
)

// A *very* simple blogging engine
func main() {
	log.Println("Starting Application")

	// Accept a command line flag "-reset"
	// This flag allows you to change the key needed to edit/delete posts
	reset := flag.Bool("reset", false, "reset the security key used to edit/delete\n")
	// Accept a command line flag "-addr :8080"
	// This flag tells the server the address to listen on
	addr := flag.String("addr", "localhost:8080", "the address/port to listen on \nuse :<port> to listen on all addresses\n")
	// Accept a command line flag "-templates ./templates"
	// This flag tells the server the path to the templates folder
	tmplPath := flag.String("tmpl", "templates", "the path to the templates folder \nfound in the src repository\n")
	// Accept a command line flag "-db ./data.db"
	// This flag tells the server the path to the sqlite database file
	dbFile := flag.String("db", "data.db", "the path to the sqlite database file \nit will be created if it does not already exist\n")

	flag.Parse()

	// Connect to the sqlite database and make sure the schema exists
	db, err := initDb(*dbFile)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally clear out the old edit key and ask the user for a new one
	if *reset {
		DeleteHashCommand(db)
	}
	EnsureHashExists(db)

	// Parse any html templates used by the application
	tmpl, err := parseTemplates(*tmplPath)
	if err != nil {
		log.Fatal(err)
	}

	// Register all of the routing handlers
	http.Handle("/", registerRoutes(db, tmpl))

	// Start the application server
	log.Println("Listening on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

	log.Println("Application Finished")
}
