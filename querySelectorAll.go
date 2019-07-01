package main

import (
	"golang.org/x/net/html"
	"regexp"
)

func querySelectorAll(root *html.Node, selector string) (nodeList []*html.Node) {
	var test func(*html.Node) bool
	switch string(selector[0]) {
	case "#":
		test = func(node *html.Node) bool {
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
		test = func(node *html.Node) bool {
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
		test = func(node *html.Node) bool {
			if node.Type == html.ElementNode && node.Data == selector {
				return true
			}
			return false
		}
	}

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