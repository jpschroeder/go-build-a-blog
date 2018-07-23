package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello World")

	data := Data{}
	data.init()
	//slug := data.create(Page{Title: "Blah bloo asdf", Date: time.Now(), Show: true, Body: "test body"})
	//fmt.Println(slug)

	slug := data.update("blah-bloo", &Page{Title: "Blah bloo", Date: time.Now(), Show: true, Body: "test body updated"})
	fmt.Println(slug)

	p := data.view(slug)
	fmt.Printf("%v", p)

	list := data.list()
	fmt.Printf("%v", list)

	fmt.Println("Done")
}
