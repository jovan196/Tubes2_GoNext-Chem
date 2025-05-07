package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeElement() {
	url := "https://little-alchemy.fandom.com/wiki/Element_Combinations"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	graph := map[string][][2]string{}

	doc.Find("table.wikitable").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() >= 2 {
				product := strings.TrimSpace(cells.Eq(0).Text())
				rawCombo := strings.TrimSpace(cells.Eq(1).Text())

				combos := strings.Split(rawCombo, "\n")
				for _, combo := range combos {
					parts := strings.Split(combo, "+")
					if len(parts) == 2 {
						a := strings.TrimSpace(parts[0])
						b := strings.TrimSpace(parts[1])
						graph[product] = append(graph[product], [2]string{a, b})
					}
				}
			}
		})
	})

	jsonFile, _ := os.Create("elements_graph.json")
	defer jsonFile.Close()
	json.NewEncoder(jsonFile).Encode(graph)

	csvFile, _ := os.Create("elements_graph.csv")
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Product", "Ingredient1", "Ingredient2"})
	for prod, pairs := range graph {
		for _, pair := range pairs {
			writer.Write([]string{prod, pair[0], pair[1]})
		}
	}

	fmt.Printf("Sukses scraping %d elemen\n", len(graph))
}
