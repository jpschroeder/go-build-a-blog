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
	// Accept a command line flag "-db ./go-build-a-blog.db"
	// This flag tells the server the path to the sqlite database file
	dbFile := flag.String("db", "go-build-a-blog.db", "the path to the sqlite database file \nit will be created if it does not already exist\n")

	flag.Parse()

	// Connect to the sqlite database and make sure the schema exists
	log.Println("Initialize Database")
	db, err := initDb(*dbFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ensure Default Data Exists")
	EnsureDefaultBlogExists(db)

	// Optionally clear out the old default key and ask the user for a new one
	if *reset {
		ResetDefaultBlogKey(db)
	}

	// Parse any html templates used by the application
	log.Println("Parse Templates")
	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatal(err)
	}

	go ExpireSessionsJob(db)

	// Register all of the routing handlers
	log.Println("Register Routes")
	http.Handle("/", registerRoutes(db, tmpl))

	// Start the application server
	log.Println("Listening on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

	log.Println("Application Finished")
}
