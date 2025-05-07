package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func ScrapeElement() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	graph := map[string][][2]string{}

	c := colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"))
	tableIndex := 0

	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++

		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			product := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if product == "" || product == "Time" || product == "Ruins" || product == "Archeologist" {
				return
			}

			h.ForEach("td:nth-of-type(2) li", func(_ int, li *colly.HTMLElement) {
				aTags := li.DOM.Find("a")
				if aTags.Length() < 4 {
					return
				}

				ingredient1 := strings.TrimSpace(aTags.Eq(1).Text())
				ingredient2 := strings.TrimSpace(aTags.Eq(3).Text())

				if ingredient1 == "" || ingredient2 == "" || ingredient1 == "Time" || ingredient2 == "Time" ||
					ingredient1 == "Ruins" || ingredient2 == "Ruins" || ingredient1 == "Archeologist" || ingredient2 == "Archeologist" {
					return
				}

				graph[product] = append(graph[product], [2]string{ingredient1, ingredient2})
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// Simpan ke JSON
	jsonFile, err := os.Create("elements_graph.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	json.NewEncoder(jsonFile).Encode(graph)

	// Simpan ke CSV
	csvFile, err := os.Create("elements_graph.csv")
	if err != nil {
		log.Fatal(err)
	}
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
