package main

//go:generate go-bindata -ignore=src templates/... static/...

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// A *very* simple blogging engine
func main() {
	log.Println("Starting Application")

	// Accept a command line flag "-db ./go-build-a-blog.db"
	// This flag tells the server the path to the sqlite database file
	dbFile := flag.String("db", "go-build-a-blog.db",
		"the path to the sqlite database file \n"+
			"it will be created if it does not already exist\n")

	// Accept a command line flag "-reset"
	// This flag allows you to change the key needed to edit/delete posts
	reset := flag.Bool("reset", false,
		"reset the security key used to edit/delete\n")

	// Accept a command line flag "-httpaddr :8080"
	// This flag tells the server the http address to listen on
	httpaddr := flag.String("httpaddr", "localhost:8080",
		"the address/port to listen on for http \n"+
			"use :<port> to listen on all addresses\n")

	// Accept a command line flag "-httpsaddr :443"
	// This flag tells the server the https address to listen on
	httpsaddr := flag.String("httpsaddr", "",
		"the address/port to listen on for https \n"+
			"use :<port> to listen on all addresses\n"+
			"this should only be used when listening publicly with proper dns address configured\n"+
			"this will generate a certificate using letsencrypt\n"+
			"the server will also listen on the -httpaddr but will redirect to https\n")

	// Accept a command line flag "-httpsaddr :443"
	// This flag tells the server the https address to listen on
	httpsdomain := flag.String("httpsdomain", "",
		"the domain to use for https\n"+
			"this flag should be used in conjunction with the -httpsaddr flag\n"+
			"this should only be used when listening publicly with proper dns address configured\n")

	devmode := flag.Bool("dev", false,
		"development mode: load html temlates from the ./templates folder instead of from the bundle\n")

	log.Println("Parsing Command Line Flags")
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
	var tmpl ExecuteTemplateFunc
	if *devmode {
		tmpl, err = parseFileTemplates()
	} else {
		tmpl, err = parseBundledTemplates()
	}
	if err != nil {
		log.Fatal(err)
	}

	go ExpireSessionsJob(db)

	// Register all of the routing handlers
	log.Println("Register Routes")
	mux := registerRoutes(db, tmpl, *devmode)

	if httpsaddr == nil || *httpsaddr == "" {
		// Use HTTP

		http.Handle("/", mux)

		// Start the application server
		log.Println("Listening on http", *httpaddr)
		log.Fatal(http.ListenAndServe(*httpaddr, nil))
	} else {
		// Use HTTPS

		log.Println("Configure Certificate")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			//Cache:  autocert.DirCache("certs"),
			Cache: SqliteCache{db: db},
		}

		if httpsdomain != nil && *httpsdomain != "" {
			certManager.HostPolicy = autocert.HostWhitelist(*httpsdomain)
		}

		server := &http.Server{
			Addr:    *httpsaddr,
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go func() {
			log.Println("Listening on http", *httpaddr)
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(*httpaddr, h))
		}()
		log.Println("Listening on https", *httpsaddr)
		log.Fatal(server.ListenAndServeTLS("", ""))
	}

	log.Println("Application Finished")
}
