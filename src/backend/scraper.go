package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Recipe represents a combination: result ← input elements
type Recipe struct {
	Result string   `json:"result"`
	Input  []string `json:"input"`
}

func ScrapeElement() {
	// collector with timeout
	c := colly.NewCollector(
		colly.UserAgent("LA2-Scraper/1.0"),
	)
	c.SetRequestTimeout(30 * time.Second)

	var recipes []Recipe
	plus := regexp.MustCompile(`\s*\+\s*`)

	// for each row in the table
	c.OnHTML("table.list-table tbody tr", func(e *colly.HTMLElement) {
		// get the product name
		result := strings.TrimSpace(e.ChildText("td:nth-child(1) a"))
		if result == "" {
			return
		}
		// iterate each <li> in the second column
		e.ForEach("td:nth-child(2) li", func(_ int, li *colly.HTMLElement) {
			// try to get all anchors to wiki pages
			elements := li.ChildAttrs("a[href^='/wiki/']", "href")
			comps := make([]string, 0, len(elements))
			for _, href := range elements {
				comps = append(comps, li.DOM.Find("a[href='"+href+"']").Text())
			}
			// fallback: split on plus if fewer than 2 anchors
			if len(comps) < 2 {
				for _, part := range plus.Split(li.Text, -1) {
					if p := strings.TrimSpace(part); p != "" {
						comps = append(comps, p)
					}
				}
			}
			if len(comps) < 2 {
				return
			}
			recipes = append(recipes, Recipe{
				Result: result,
				Input:  uniqueSort(comps),
			})
		})
	})

	// visit the page
	if err := c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"); err != nil {
		log.Fatalf("Failed to visit page: %v", err)
	}

	// sort recipes by result, then inputs
	sort.Slice(recipes, func(i, j int) bool {
		if recipes[i].Result != recipes[j].Result {
			return recipes[i].Result < recipes[j].Result
		}
		return strings.Join(recipes[i].Input, ",") < strings.Join(recipes[j].Input, ",")
	})

	// write output JSON
	outFile := "elements_graph.json"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", outFile, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(map[string]interface{}{"recipes": recipes}); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}

	fmt.Printf("✅ Parsed %d recipes → %s\n", len(recipes), outFile)
}

// uniqueSort dedups and sorts a slice of strings
func uniqueSort(in []string) []string {
	set := make(map[string]struct{}, len(in))
	for _, s := range in {
		set[s] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
