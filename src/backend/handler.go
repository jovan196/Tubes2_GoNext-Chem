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

// writeJSONResponse menulis v sebagai JSON ke ResponseWriter
func writeJSONResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func buildStepsFromTrace(node *TraceNode) []RecipeStep {
	visited := make(map[string]bool)
	var steps []RecipeStep
	var dfs func(*TraceNode)
	dfs = func(n *TraceNode) {
		if n == nil || visited[n.Product] {
			return
		}
		visited[n.Product] = true
		// jika leaf (basic), parent mungkin nil
		if n.Parent[0] == nil || n.Parent[1] == nil {
			return
		}
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
			var out []SearchResponse
			for _, n := range nodes {
				out = append(out, SearchResponse{Result: req.Target, Steps: buildStepsFromTrace(n)})
			}
			writeJSONResponse(w, out)
			return
		}
		if node := BFS(req.Target); node != nil {
			writeJSONResponse(w, []SearchResponse{{Result: req.Target, Steps: buildStepsFromTrace(node)}})
			return
		}
		http.Error(w, "Element not reachable", http.StatusNotFound)

	case "dfs":
		if req.Mode == "multiple" {
			nodes := MultiDFS_Trace(req.Target, req.Max)
			var out []SearchResponse
			for _, n := range nodes {
				out = append(out, SearchResponse{Result: req.Target, Steps: buildStepsFromTrace(n)})
			}
			writeJSONResponse(w, out)
			return
		}
		if node := DFS(req.Target); node != nil {
			writeJSONResponse(w, []SearchResponse{{Result: req.Target, Steps: buildStepsFromTrace(node)}})
			return
		}
		http.Error(w, "Element not reachable", http.StatusNotFound)

	default:
		http.Error(w, "Only BFS and DFS are implemented", http.StatusNotImplemented)
	}
}
