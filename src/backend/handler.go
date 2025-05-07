package main

import (
	"encoding/json"
	"net/http"
)

type SearchRequest struct {
	Target    string `json:"target"`
	Algorithm string `json:"algorithm"` // BFS / DFS
	Mode      string `json:"mode"`      // single / multiple
	Max       int    `json:"max"`       // only used in multiple
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

	// Dummy logic, nanti ganti dengan BFS/DFS
	response := []SearchResponse{
		{
			Result: req.Target,
			Steps: [][]string{
				{"Mud", "Fire"},
				{"Clay", "Stone"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
