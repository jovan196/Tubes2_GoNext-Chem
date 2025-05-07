package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SearchRequest struct {
	Target    string `json:"target"`
	Algorithm string `json:"algorithm"`
	Mode      string `json:"mode"`
	Max       int    `json:"max"`
}

type SearchResponse struct {
	Result string     `json:"result"`
	Steps  [][]string `json:"steps"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	algo := strings.ToLower(req.Algorithm)
	if req.Mode == "multiple" {
		results := MultiBFS(req.Target, req.Max)
		var responses []SearchResponse
		for _, path := range results {
			steps := [][]string{}
			for _, pair := range path {
				steps = append(steps, []string{pair[0], pair[1]})
			}
			responses = append(responses, SearchResponse{
				Result: req.Target,
				Steps:  steps,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
		return
	}

	if algo == "bfs" {
		path, found := BFS(req.Target)
		if !found {
			http.Error(w, "Element not reachable", http.StatusNotFound)
			return
		}

		steps := [][]string{}
		for _, pair := range path {
			steps = append(steps, []string{pair[0], pair[1]})
		}

		response := []SearchResponse{
			{
				Result: req.Target,
				Steps:  steps,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Only BFS is implemented", http.StatusNotImplemented)
	}
}
