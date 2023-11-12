package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

type WebScraper struct {
	baseURL string
	*http.Client
	Headers map[string]string
}

func (c *WebScraper) Visit(endpoint string) ([]*html.Node, int, error) {
	fmt.Println("Scrap from: ", c.baseURL+endpoint)
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Do request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	defer resp.Body.Close()

	// Read and parse html body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	nodes, err := Parse(string(body))
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return nodes, resp.StatusCode, nil
}

func (c *WebScraper) SetHeader(key, value string) {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers[key] = value
	return
}

func (c *WebScraper) RemoveHeader(key string) {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	delete(c.Headers, key)
	return
}

func Parse(text string) ([]*html.Node, error) {
	nodes := make([]*html.Node, 1)
	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			nodes = append(nodes, n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nodes, nil
}

func GetByNameTag(nodes []*html.Node, tag string) ([]*html.Node, error) {
	newNodes := make([]*html.Node, 0)
	for _, node := range nodes {
		if node != nil {
			if node.DataAtom.String() == tag {
				newNodes = append(newNodes, node)
			}
		}
	}

	return newNodes, nil
}

func GetByAttribute(nodes []*html.Node, attribute html.Attribute) ([]*html.Node, error) {
	newNodes := make([]*html.Node, 0)
	for _, node := range nodes {
		if node != nil {
			for _, attr := range node.Attr {
				if strings.Contains(attr.Key, attribute.Key) && strings.Contains(attr.Val, attribute.Val) {
					newNodes = append(newNodes, node)
				}
			}
		}
	}

	return newNodes, nil
}
