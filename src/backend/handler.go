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
			result := SearchResponse{
				Result:       req.Target,
				Tree:         root,
				TimeMs:       time.Since(start).Milliseconds(),
				VisitedCount: LastBFSVisited,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]SearchResponse{result})
			return
		}

		node := BFS(req.Target)
		var tree *OutputNode
		var visitCount int
		if node == nil {
			// tree = MultiBFS_Trace(req.Target, 1)
			// if tree == nil {
			http.Error(w, createError("Element "+req.Target+" not reachable", http.StatusNotFound), http.StatusNotFound)
			return
			//}
			visitCount = LastBFSVisited
		} else {
			tree = toOutputTree(node)
			visitCount = LastBFSVisited
		}

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
		http.Error(w, createError("Only BFS implemented (revised)", http.StatusNotImplemented), http.StatusNotImplemented)
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
