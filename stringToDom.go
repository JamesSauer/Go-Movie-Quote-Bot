package main

import (
	"golang.org/x/net/html"
	"log"
	"strings"
)

func stringToDom(str string) (document *html.Node) {
	document, err := html.Parse(strings.NewReader(str))
	if err != nil {
		log.Fatal(err)
	}
	return
}