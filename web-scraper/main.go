package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	webScraper := WebScraper{baseURL: "https://www.otodom.pl/", Client: &http.Client{}}

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

	flats := make([]Flat, 0)
	flatCh := make(chan Flat)
	for i := 1; i <= page; i++ {
		nodes, _, err = webScraper.Visit("pl/wyniki/sprzedaz/mieszkanie/slaskie/katowice/katowice/katowice?viewType=listing&limit=72&page=" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		nodes, err = GetByAttribute(nodes, html.Attribute{Namespace: "", Key: "data-cy", Val: "listing-item-link"})
		if err != nil {
			panic(err)
		}

		// Get specific offer
		for _, node := range nodes {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if node != nil && len(node.Attr) >= 2 {
						go fetchOffer(&webScraper, flatCh, node.Attr[1].Val)
					}
				}
			}
		}
	}

	// CLOSE CHANNEL
	close(flatCh)

	// CONSUME FLATS
	for flat := range flatCh {
		flats = append(flats, flat)
	}

	//Save flats to csv
	file, err := os.Create("flats.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	writer := csv.NewWriter(file)

	if err := writer.Write([]string{"Nazwa", "Powierzchnia", "Liczba pokoi"}); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, flat := range flats {
		row := []string{flat.Name, fmt.Sprintf("%f", flat.Surface), strconv.Itoa(flat.NoOfRooms)}
		if err := writer.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}

type Flat struct {
	Name      string  `json:"name"`
	Surface   float64 `json:"surface"`
	NoOfRooms int     `json:"no_of_rooms"`
}

func fetchOffer(webScraper *WebScraper, flatCh chan<- Flat, url string) {
	offerNodes, status, err := webScraper.Visit(url)
	if err != nil {
		panic(err)
	}
	if status == http.StatusOK {
		// Make a CSV file
		n, err := GetByAttribute(offerNodes, html.Attribute{Namespace: "", Key: "data-testid", Val: "table-value-area"})
		if err != nil {
			panic(err)
		}

		surface, err := strconv.ParseFloat(strings.Replace(n[0].FirstChild.Data[:len(n[0].FirstChild.Data)-4], ",", ".", 1), 64)
		if err != nil {
			surface = 0.0
		}

		n, err = GetByAttribute(offerNodes, html.Attribute{Namespace: "", Key: "data-testid", Val: "table-value-rooms_num"})
		if err != nil {
			panic(err)
		}

		noOfRooms, err := strconv.Atoi(strings.Trim(n[0].LastChild.Data, " "))
		if err != nil {
			noOfRooms, err = strconv.Atoi(strings.Trim(n[0].LastChild.FirstChild.Data, " "))
			if err != nil {
				noOfRooms = 0
			}
		}

		flatCh <- Flat{
			Name:      url,
			Surface:   surface,
			NoOfRooms: noOfRooms,
		}
	}
}
