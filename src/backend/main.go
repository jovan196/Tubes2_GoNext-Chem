package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Scrape & simpan graph (jika perlu setiap startup)
	ScrapeElement()

	// Load ke memori
	if err := LoadGraph("elements_graph.json"); err != nil {
		log.Fatal("LoadGraph:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", SearchHandler)

	// izinkan CORS untuk frontend Next.js
	handler := cors.Default().Handler(mux)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
