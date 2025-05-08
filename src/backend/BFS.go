package main

import (
	"sync"
)

type TraceNode struct {
	Product string
	From    [2]string
	Parent  [2]*TraceNode
}

func BFS(target string) *TraceNode {
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

func MultiBFS_Trace(target string, max int) []*TraceNode {
	var results []*TraceNode
	var mu sync.Mutex
	var wg sync.WaitGroup
	resultChan := make(chan *TraceNode, max)

	visitedGlobal := map[string]bool{
		"Air": true, "Water": true, "Earth": true, "Fire": true, "Time": true,
	}

	for product, recipes := range Graph {
		if product == target {
			for _, pair := range recipes {
				wg.Add(1)
				go func(pair [2]string) {
					defer wg.Done()

					queue := []*TraceNode{
						{Product: pair[0]},
						{Product: pair[1]},
					}

					visited := make(map[string]bool)
					for k, v := range visitedGlobal {
						visited[k] = v
					}
					visited[pair[0]] = true
					visited[pair[1]] = true

					for len(queue) > 0 {
						curr := queue[0]
						queue = queue[1:]

						for prod, recs := range Graph {
							for _, rec := range recs {
								if (curr.Product == rec[0] && visited[rec[1]]) || (curr.Product == rec[1] && visited[rec[0]]) {
									node := &TraceNode{
										Product: product,
										From:    [2]string{pair[0], pair[1]},
										Parent:  [2]*TraceNode{findNode(queue, pair[0]), findNode(queue, pair[1])},
									}

									if prod == target {
										mu.Lock()
										if len(results) < max {
											results = append(results, node)
											resultChan <- node
										}
										mu.Unlock()
										return
									}
									if !visited[prod] {
										visited[prod] = true
										queue = append(queue, node)
									}
								}
							}
						}
					}
				}(pair)
			}
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for node := range resultChan {
		results = append(results, node)
		if len(results) >= max {
			break
		}
	}

	return results
}

func findNode(queue []*TraceNode, name string) *TraceNode {
	for i := len(queue) - 1; i >= 0; i-- {
		if queue[i].Product == name {
			return queue[i]
		}
	}
	return nil
}
