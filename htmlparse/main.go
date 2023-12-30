package main

import (
	"flag"
	"fmt"
	"os"
	"part3/link"
)

func main() {
	filename := flag.String("file", "ex1.html", "The HTML file to parse links from")
	flag.Parse() //file名を指定してパース可能
	s, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	links, err := link.Parse(s)
	fmt.Println(links)
	if err != nil {
		panic(err)
	}

	for _, i := range links {
		fmt.Println("Href: ", i.Href)
		fmt.Println("Text: ", i.Text)
	}
	return
}
