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

type RecipeStep struct {
	Product     string   `json:"product"`
	Ingredients []string `json:"ingredients"`
}

type SearchResponse struct {
	Result string       `json:"result"`
	Steps  []RecipeStep `json:"steps"`
}

func buildStepsFromTrace(node *TraceNode) []RecipeStep {
	visited := map[string]bool{}
	var steps []RecipeStep

	var dfs func(*TraceNode)
	dfs = func(n *TraceNode) {
		if n == nil || visited[n.Product] || n.Parent[0] == nil || n.Parent[1] == nil {
			return
		}
		visited[n.Product] = true
		dfs(n.Parent[0])
		dfs(n.Parent[1])
		steps = append(steps, RecipeStep{
			Product:     n.Product,
			Ingredients: []string{n.From[0], n.From[1]},
		})
	}
	dfs(node)

	return steps
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
	if algo != "bfs" {
		http.Error(w, "Only BFS is implemented", http.StatusNotImplemented)
		return
	}

	if req.Mode == "multiple" {
		results := MultiBFS_Trace(req.Target, req.Max) // new MultiBFS version returning []*TraceNode
		var responses []SearchResponse
		for _, node := range results {
			responses = append(responses, SearchResponse{
				Result: req.Target,
				Steps:  buildStepsFromTrace(node),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
		return
	}

	// mode == "single"
	node := BFS(req.Target)
	if node == nil {
		http.Error(w, "Element not reachable", http.StatusNotFound)
		return
	}

	response := SearchResponse{
		Result: req.Target,
		Steps:  buildStepsFromTrace(node),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]SearchResponse{response})
}
