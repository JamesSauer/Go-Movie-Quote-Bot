package main

import (
	"errors"
	"golang.org/x/net/html"
	"regexp"
)


/* TODO: Look into using goquery to replace the functions in this file. https://github.com/PuerkitoBio/goquery
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

// Strips the tags off the text within a given node.
// Without this, inline tags like <b>...</b> would screw up the quotes.
func extractText(root *html.Node) (text string) {
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(root)
	return
}

// Extracts the "href" attribute of the first a-tag child it finds within a given container node.
// This function assumes every node it checks to have only one child.
// Siblings do not get checked.
func extractLink(node *html.Node) (string, error) {
	// TODO: Make use of querySelector for this, as it does most of what this function does as well.
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				findRedlink := regexp.MustCompile("redlink=1")
				if findRedlink.MatchString(a.Val) {
					return "", errors.New("Found link was a red link")
				}
				return a.Val, nil
			}
		}
	}
	if node.FirstChild == nil {
		return "", errors.New("No link found")
	}
	return extractLink(node.FirstChild)
}

// Because html.Node.NextSibling doesn't differentiate between text and element nodes.
func getNextElementSibling(node *html.Node) (sibling *html.Node) {
	for sibling = node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type == html.ElementNode {
			return
		}
	}
	return nil
}

func getFirstElementChild(node *html.Node) (firstChild *html.Node) {
	// TODO: This function does almost the exact same thing as the one above. Merge them!
	for firstChild = node.FirstChild; firstChild != nil; firstChild = firstChild.NextSibling {
		if firstChild.Type == html.ElementNode {
			return
		}
	}
	return nil
}