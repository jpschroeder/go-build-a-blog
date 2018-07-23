package main

import (
	"fmt"
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

	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		//slug := data.create(Page{Title: "Blah bloo asdf", Date: time.Now(), Show: true, Body: "test body"})
		//fmt.Println(slug)

		slug := data.update("blah-bloo", &Page{Title: "Blah bloo", Date: time.Now(), Show: true, Body: "test body updated"})
		fmt.Println(slug)

		p := data.view(slug)
		fmt.Printf("%v", p)

		list := data.list()
		fmt.Printf("%v", list)
	*/

	fmt.Println("Done")
}
