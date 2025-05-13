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

type SearchResponse struct {
	Result       string      `json:"result"`
	Tree         *OutputNode `json:"tree"` // untuk multiple bisa []*OutputNode
	TimeMs       int64       `json:"timeMs"`
	VisitedCount int         `json:"visitedCount"`
}

// countOutputNodes recursively counts nodes in the output tree
func countOutputNodes(n *OutputNode) int {
	if n == nil {
		return 0
	}
	count := 1
	for _, c := range n.Children {
		count += countOutputNodes(c)
	}
	return count
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
	start := time.Now()

	switch algo {
	case "bfs":
		if req.Mode == "multiple" {
			root := MultiBFS_Trace(req.Target, req.Max)
			visitedCount := countOutputNodes(root)
			// Uji waktu
			// time.Sleep(1 * time.Second)
			result := SearchResponse{
				Result:       req.Target,
				Tree:         root,
				TimeMs:       time.Since(start).Milliseconds(),
				VisitedCount: visitedCount,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]SearchResponse{result})
			return
		}

		node := BFS(req.Target)
		var tree *OutputNode
		var visitCount int
		if node == nil {
			tree = nil
			visitCount = 0

		} else {
			tree = toOutputTree(node)
			visitCount = LastBFSVisited
		}
		// Uji waktu
		// time.Sleep(1 * time.Second)
		resp := SearchResponse{
			Result:       req.Target,
			Tree:         tree,
			TimeMs:       time.Since(start).Milliseconds(),
			VisitedCount: visitCount,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{resp})
		return

	case "dfs":
		if req.Mode == "multiple" {
			results := MultiDFS_Trace(req.Target, req.Max)
			var output *OutputNode
			var visitedCount int
			if len(results) == 0 {
				output = nil
				visitedCount = 0
			} else {
				output = mergeTraceTrees(results)
				visitedCount = countOutputNodes(output)
			}
			// Uji waktu
			// time.Sleep(1 * time.Second)
			result := SearchResponse{
				Result:       req.Target,
				Tree:         output,
				TimeMs:       time.Since(start).Milliseconds(),
				VisitedCount: visitedCount,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]SearchResponse{result})
			return
		}

		node := DFS(req.Target) // Ensure this is inside the DFS case block
		var tree *OutputNode
		var visitCount int // Ensure this is inside the DFS case block

		if node == nil {
			tree = nil
			visitCount = 0
		} else {
			tree = toOutputTree(node)
			visitCount = LastDFSVisited
		}
		// Uji waktu
		// time.Sleep(1 * time.Second)
		resp := SearchResponse{
			Result:       req.Target,
			Tree:         tree,
			TimeMs:       time.Since(start).Milliseconds(),
			VisitedCount: visitCount,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResponse{resp})
		return

	default:
		http.Error(w, createError("Only BFS and DFS implemented", http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func createError(message string, statusCode int) string {
	response := map[string]interface{}{
		"error":      true,
		"message":    message,
		"statusCode": statusCode,
	}
	errorResponse, _ := json.Marshal(response)
	return string(errorResponse)
}
