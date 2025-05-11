package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	ScrapeElement()

	var err error

	err = LoadGraph("elements_graph.json")
	if err != nil {
		log.Fatal("Failed to load graph:", err)
	}

	err = LoadTier("elements_tier.json")
	if err != nil {
		log.Fatal("Failed to load tier:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", SearchHandler)

	handler := cors.Default().Handler(mux)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
