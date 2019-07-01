package main

import (
	"golang.org/x/net/html"
	"regexp"
)


/* TODO: Look into using goquery to replace the functions below. https://github.com/PuerkitoBio/goquery
*/

// A helper for use in the querySelector functions.
// It returns a function that tests for either a specific id, class or tag name depending on selector.
func getTest(selector string) func(*html.Node) bool {
	// TODO: The regexps in here aren't very robust. Improve them!
	// TODO: Make it accept compound selectors, like ".mw-headline.highlightes" or "a.outbound".
	switch string(selector[0]) {
	case "#":
		// Select by id.
		return func(node *html.Node) bool {
			if node.Type == html.ElementNode {
				for _, attr := range node.Attr {
					if attr.Key == "id" && attr.Val == selector[1:] {
						return true
					}
					return false
				}
			}
			return false
		}
	case ".":
		// Select by class.
		return func(node *html.Node) bool {
			if node.Type == html.ElementNode {
				for _, attr := range node.Attr {
					if attr.Key == "class" {
						matches, _ := regexp.MatchString(selector[1:], attr.Val)
						if matches {
							return true
						}
						return false
					}
				}
			}
			return false
		}
	default:
		// Select by tag name.
		return func(node *html.Node) bool {
			if node.Type == html.ElementNode && node.Data == selector {
				return true
			}
			return false
		}
	}
}


func querySelectorAll(root *html.Node, selector string) (nodeList []*html.Node) {
	test := getTest(selector)

	nodeList = make([]*html.Node, 0)

	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if test(node) {
			nodeList = append(nodeList, node)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(root)

	return
}