package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
)

func main() {
	webScraper := &WebScraper{baseURL: "https://www.otodom.pl/", Client: &http.Client{}}

	// Set Agent Header
	webScraper.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")

	// Scrap amount of pages
	nodes, _, err := webScraper.Visit("pl/wyniki/sprzedaz/mieszkanie/slaskie/katowice/katowice/katowice?viewType=listing&limit=72")

	if err != nil {
		panic(err)
	}

	nodes, err = GetByNameTag(nodes, "a")
	if err != nil {
		panic(err)
	}

	nodes, err = GetByAttribute(nodes, html.Attribute{Namespace: "", Key: "data-cy", Val: "pagination"})
	if err != nil {
		panic(err)
	}

	page, err := strconv.Atoi(nodes[len(nodes)-1].FirstChild.Data)

	fmt.Println("Number of pages:", strconv.Itoa(page))
}
