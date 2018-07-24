package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Application")

	data := Data{}
	data.init()
	handlers := Handlers{data: data}
	handlers.init()
	http.Handle("/", handlers.registerRoutes())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Application Finished")
}
