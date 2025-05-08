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
		if n == nil || visited[n.Product] {
			return
		}
		// leaf nodes (basic) may have nil parents
		if n.Parent[0] == nil || n.Parent[1] == nil {
			visited[n.Product] = true
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	algo := strings.ToLower(req.Algorithm)
	switch algo {
	case "bfs":
		if req.Mode == "multiple" {
			nodes := MultiBFS_Trace(req.Target, req.Max)
			var resp []SearchResponse
			for _, node := range nodes {
				resp = append(resp, SearchResponse{
					Result: req.Target,
					Steps:  buildStepsFromTrace(node),
				})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		// single
		node := BFS(req.Target)
		if node == nil {
			http.Error(w, "Element not reachable", http.StatusNotFound)
			return
		}
		result := SearchResponse{Result: req.Target, Steps: buildStepsFromTrace(node)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{result})
		return

	case "dfs":
		if req.Mode == "multiple" {
			nodes := MultiDFS_Trace(req.Target, req.Max)
			var resp []SearchResponse
			for _, node := range nodes {
				resp = append(resp, SearchResponse{
					Result: req.Target,
					Steps:  buildStepsFromTrace(node),
				})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		node := DFS(req.Target)
		if node == nil {
			http.Error(w, "Element not reachable", http.StatusNotFound)
			return
		}
		result := SearchResponse{Result: req.Target, Steps: buildStepsFromTrace(node)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{result})
		return

	default:
		http.Error(w, "Only BFS and DFS are implemented", http.StatusNotImplemented)
		return
	}
}
