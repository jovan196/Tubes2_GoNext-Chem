package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var LastBFSVisited int

type TraceNode struct {
	Product string
	From    [2]string
	Parent  [2]*TraceNode
	Depth   int
}

func BFS(target string) *TraceNode {
	LastBFSVisited = 0
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
				if Tier[a] >= Tier[target] && Tier[b] >= Tier[target] {
					continue
				}
				if (curr.Product == a && visited[b]) || (curr.Product == b && visited[a]) {
					left := allNodes[a]
					right := allNodes[b]
					if left == nil || right == nil {
						continue
					}

					node := &TraceNode{
						Product: product,
						From:    [2]string{a, b},
						Parent:  [2]*TraceNode{left, right},
						Depth:   1 + max(left.Depth, right.Depth),
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

func MultiBFS_Trace(target string, maxResults int) []*TraceNode {
	// Reset counter for this search operation
	LastBFSVisited = 0

	basic := []string{"Air", "Water", "Earth", "Fire", "Time"}
	nodes := make(map[string][]*TraceNode)
	queue := [][]*TraceNode{}
	results := []*TraceNode{}
	seenHash := make(map[string]bool)
	var counter int32

	for _, b := range basic {
		n := &TraceNode{Product: b}
		nodes[b] = []*TraceNode{n}
		queue = append(queue, []*TraceNode{n})
	}

	var mu sync.Mutex
	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	for len(queue) > 0 && int(atomic.LoadInt32(&counter)) < maxResults {
		currLevel := queue[0]
		queue = queue[1:]

		// Increment the visit counter for each level
		atomic.AddInt32(&counter, 1)
		LastBFSVisited++
		nextLevel := []*TraceNode{}

		for _, curr := range currLevel {
			for outProd, recs := range Graph {
				for _, pair := range recs {
					a, b := pair[0], pair[1]
					if curr.Product != a && curr.Product != b {
						continue
					}

					mu.Lock()
					listA := nodes[a]
					listB := nodes[b]
					mu.Unlock()

					if len(listA) == 0 || len(listB) == 0 {
						continue
					}

					for _, ta := range listA {
						for _, tb := range listB {
							if Tier[a] >= Tier[target] || Tier[b] >= Tier[target] {
								continue
							}
							sem <- struct{}{}
							wg.Add(1)
							go func(a, b, out string, ta, tb *TraceNode) {
								defer wg.Done()
								defer func() { <-sem }()

								if containProduct(ta, out) || containProduct(tb, out) {
									return
								}

								mu.Lock()
								node := &TraceNode{
									Product: out,
									From:    [2]string{a, b},
									Parent:  [2]*TraceNode{ta, tb},
									Depth:   1 + max(ta.Depth, tb.Depth),
								}
								hash := hashSubtree(node)
								if seenHash[hash] {
									mu.Unlock()
									return
								}
								seenHash[hash] = true
								nodes[out] = append(nodes[out], node)
								if out == target && int(atomic.LoadInt32(&counter)) < maxResults {
									results = append(results, node)
									atomic.AddInt32(&counter, 1)
								}
								mu.Unlock()
								nextLevel = append(nextLevel, node)
							}(a, b, outProd, ta, tb)
						}
					}
				}
			}
		}

		wg.Wait()
		if len(nextLevel) > 0 {
			queue = append(queue, nextLevel)
		}
	}

	return results
}

func hashSubtree(n *TraceNode) string {
	if n == nil {
		return ""
	}
	if n.Parent[0] == nil && n.Parent[1] == nil {
		return n.Product
	}
	l := hashSubtree(n.Parent[0])
	r := hashSubtree(n.Parent[1])
	if l > r {
		l, r = r, l
	}
	return fmt.Sprintf("%s(%s+%s)", n.Product, l, r)
}

func containProduct(n *TraceNode, prod string) bool {
	if n == nil {
		return false
	}
	if n.Product == prod {
		return true
	}
	return containProduct(n.Parent[0], prod) ||
		containProduct(n.Parent[1], prod)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
