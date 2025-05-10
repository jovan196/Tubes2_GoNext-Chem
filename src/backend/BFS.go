package main

import (
	"fmt"
)

var LastBFSVisited int

type TraceNode struct {
	Product string
	From    [2]string
	Parent  [2]*TraceNode
}

func BFS(target string) *TraceNode {
	LastBFSVisited = 0
	// Inisialisasi node dasar
	basicElements := []string{"Air", "Water", "Earth", "Fire", "Time"}
	allNodes := map[string]*TraceNode{}
	queue := []*TraceNode{}

	for _, elem := range basicElements {
		node := &TraceNode{Product: elem}
		allNodes[elem] = node
		queue = append(queue, node)
	}

	visited := map[string]bool{}
	for _, elem := range basicElements {
		visited[elem] = true
	}

	for len(queue) > 0 {
		LastBFSVisited++
		curr := queue[0]
		queue = queue[1:]

		for product, recipes := range Graph {
			if visited[product] {
				continue
			}
			for _, pair := range recipes {
				a, b := pair[0], pair[1]
				if (curr.Product == a && visited[b]) || (curr.Product == b && visited[a]) {
					// Cek apakah bahan-bahan sudah pernah dibuat
					left := allNodes[a]
					right := allNodes[b]
					if left == nil || right == nil {
						continue
					}

					node := &TraceNode{
						Product: product,
						From:    [2]string{a, b},
						Parent:  [2]*TraceNode{left, right},
					}

					allNodes[product] = node
					if product == target {
						return node
					}

					visited[product] = true
					queue = append(queue, node)
				}
			}
		}
	}

	return nil
}

// MultiBFS_Trace mengembalikan hingga maxResults jalur unik menuju target
func MultiBFS_Trace(target string, maxResults int) []*TraceNode {
	// 1) Siapkan map product → list of trace trees
	nodesMap := make(map[string][]*TraceNode)
	// 2) Queue produk yang akan dieksplor
	queue := []string{}
	// 3) Inisialisasi elemen dasar
	basic := []string{"Air", "Water", "Earth", "Fire", "Time"}
	for _, b := range basic {
		tn := &TraceNode{Product: b}
		nodesMap[b] = []*TraceNode{tn}
		queue = append(queue, b)
	}

	// 4) Hasil akhir
	var results []*TraceNode

	// 5) Untuk menghindari duplikasi sama persis (a, b, dan sub-tree yang sama)
	visitedComb := make(map[string]map[string]bool) // outProd → map[key]bool
	totalCombinations := 0
	for _, rec := range Graph[target] {
		a, b := rec[0], rec[1]
		if a == target || b == target {
			continue
		}
		totalCombinations++
	}
	if maxResults > totalCombinations {
		// jika maxResults lebih besar dari total kombinasi, set ke total kombinasi
		// maxResults = totalCombinations
	}

	// 6) BFS: level-order expand sampai kumpulkan maxResults
	for len(queue) > 0 && len(results) < maxResults {
		LastBFSVisited++
		prod := queue[0]
		queue = queue[1:]

		// Cek semua resep yang menghasilkan sesuatu dari `prod` sebagai salah satu bahan
		for outProd, recs := range Graph {
			for _, rec := range recs {
				a, b := rec[0], rec[1]
				// hanya recipes yang melibatkan prod
				if a != prod && b != prod {
					continue
				}
				tracesA, okA := nodesMap[a]
				tracesB, okB := nodesMap[b]
				// butuh keduanya punya trace
				if !okA || !okB {
					continue
				}
				// siapkan visitedComb[outProd]
				if visitedComb[outProd] == nil {
					visitedComb[outProd] = make(map[string]bool)
				}
				// buat kombinasi semua sub-tree A × sub-tree B
				for _, ta := range tracesA {
					for _, tb := range tracesB {
						// cegah self-loop apapun kedalam level deeper
						if containsProduct(ta, outProd) || containsProduct(tb, outProd) {
							continue
						}
						// unique key per kombinasi sub-tree pointer
						key := fmt.Sprintf("%s|%p|%p", outProd, ta, tb)
						if visitedComb[outProd][key] {
							continue
						}
						visitedComb[outProd][key] = true

						// bangun trace baru
						newTrace := &TraceNode{
							Product: outProd,
							From:    [2]string{a, b},
							Parent:  [2]*TraceNode{ta, tb},
						}
						// simpan
						nodesMap[outProd] = append(nodesMap[outProd], newTrace)

						// jika target, tambahkan ke results
						if outProd == target {
							results = append(results, newTrace)
							if len(results) >= maxResults {
								return results
							}
						} else {
							// enqueue outProd untuk level berikut
							queue = append(queue, outProd)
						}
					}
				}
			}
		}
	}

	return results
}
