package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// rawGraph hanya untuk menangkap top‐level "recipes"
type rawGraph struct {
	Recipes []Recipe `json:"recipes"`
}

// Graph akan diisi dengan mapping product → daftar pasangan 2‐input
var Graph map[string][][2]string

// LoadGraph membaca file JSON scraper, mengisi Graph
func LoadGraph(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// decode ke struktur per‐recipe
	var rg rawGraph
	dec := json.NewDecoder(f)
	if err := dec.Decode(&rg); err != nil {
		return err
	}

	// bangun Graph
	Graph = make(map[string][][2]string, len(rg.Recipes))
	for _, r := range rg.Recipes {
		// hanya yang inputnya tepat 2 elemen
		if len(r.Input) == 2 {
			pair := [2]string{r.Input[0], r.Input[1]}
			Graph[r.Result] = append(Graph[r.Result], pair)
		}
	}

	fmt.Printf("Loaded %d recipes into %d products from %s\n",
		len(rg.Recipes), len(Graph), filename)
	return nil
}
