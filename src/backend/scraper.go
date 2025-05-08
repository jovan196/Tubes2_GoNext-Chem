package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	// 1) Fetch page with a real User-Agent
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/115.0.0.0 Safari/537.36",
	)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Regex pasangan Elemen+Elemen
	re := regexp.MustCompile(`([A-Za-z ]+)\s*\+\s*([A-Za-z ]+)`)

	graph := make(map[string][][2]string)

	// loop tabel, cari baris yang mengandung minimal 2 kolom
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 2 {
				return
			}
			product := strings.TrimSpace(cells.Eq(0).Text())
			rawCombo := strings.TrimSpace(cells.Eq(1).Text())

			
			for _, match := range re.FindAllStringSubmatch(rawCombo, -1) {
				a := strings.TrimSpace(match[1])
				b := strings.TrimSpace(match[2])
				graph[product] = append(graph[product], [2]string{a, b})
			}
		})
	})

	// Write JSON
	jf, _ := os.Create("elements_graph.json")
	defer jf.Close()
	json.NewEncoder(jf).Encode(graph)
	w.Flush()

	fmt.Printf("Extracted %d products\n", len(graph))
}
