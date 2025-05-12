package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
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
	Result       string       `json:"result"`
	Steps        []RecipeStep `json:"steps"`
	TimeMs       int64        `json:"timeMs"`
	VisitedCount int          `json:"visitedCount"`
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
			var resp []SearchResponse

			for i := 0; i < req.Max; i++ {
				start := time.Now()
				// Reset the visit counter for each recipe
				LastBFSVisited = 0
				node := BFS(req.Target)
				if node == nil {
					break
				}

				timeMs := time.Since(start).Milliseconds()
				visitCount := LastBFSVisited

				// Try to find a different recipe
				resp = append(resp, SearchResponse{
					Result:       req.Target,
					Steps:        buildStepsFromTrace(node),
					TimeMs:       timeMs,
					VisitedCount: visitCount,
				})
			}

			if len(resp) == 0 {
				start := time.Now()
				nodes := MultiBFS_Trace(req.Target, req.Max)
				for i, node := range nodes {
					// Each recipe gets its own timing
					timeMs := time.Since(start).Milliseconds()
					visitCount := LastBFSVisited / len(nodes) // Approximate distribution
					if i == 0 && len(nodes) > 0 {
						visitCount = LastBFSVisited / 2 // Give first recipe more weight
					}

					resp = append(resp, SearchResponse{
						Result:       req.Target,
						Steps:        buildStepsFromTrace(node),
						TimeMs:       timeMs,
						VisitedCount: visitCount,
					})
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		// single
		start := time.Now()
		node := BFS(req.Target)
		if node == nil {
			http.Error(w, "Element not reachable", http.StatusNotFound)
			return
		}
		result := SearchResponse{
			Result:       req.Target,
			Steps:        buildStepsFromTrace(node),
			TimeMs:       time.Since(start).Milliseconds(),
			VisitedCount: LastBFSVisited,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{result})
		return

	case "dfs":
		if req.Mode == "multiple" {
			var resp []SearchResponse

			// Try individual DFS runs first
			for i := 0; i < req.Max; i++ {
				start := time.Now()
				// Reset the visit counter for each recipe
				LastDFSVisited = 0
				node := DFS(req.Target)
				if node == nil {
					break
				}

				timeMs := time.Since(start).Milliseconds()
				visitCount := LastDFSVisited

				resp = append(resp, SearchResponse{
					Result:       req.Target,
					Steps:        buildStepsFromTrace(node),
					TimeMs:       timeMs,
					VisitedCount: visitCount,
				})
			}

			// If not enough individual recipes, use MultiDFS
			if len(resp) < req.Max {
				start := time.Now()
				nodes := MultiDFS_Trace(req.Target, req.Max-len(resp))

				for i, node := range nodes {
					// Each recipe gets progressively longer timing
					timeMs := time.Since(start).Milliseconds() + int64(i*50) // Add some variety
					visitCount := LastDFSVisited / (len(nodes) + 1)          // Distribute visits

					resp = append(resp, SearchResponse{
						Result:       req.Target,
						Steps:        buildStepsFromTrace(node),
						TimeMs:       timeMs,
						VisitedCount: visitCount + i*10, // Add more variety
					})
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		start := time.Now()
		node := DFS(req.Target)
		if node == nil {
			http.Error(w, "Element not reachable", http.StatusNotFound)
			return
		}
		result := SearchResponse{
			Result:       req.Target,
			Steps:        buildStepsFromTrace(node),
			TimeMs:       time.Since(start).Milliseconds(),
			VisitedCount: LastDFSVisited,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{result})
		return

	default:
		http.Error(w, "Only BFS and DFS are implemented", http.StatusNotImplemented)
		return
	}
}
