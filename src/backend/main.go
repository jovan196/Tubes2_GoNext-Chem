package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", SearchHandler)

	// CORS middleware biar bisa diakses dari React
	handler := cors.Default().Handler(mux)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
