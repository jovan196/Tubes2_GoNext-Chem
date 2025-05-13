package main

import (
	"runtime"
	"sync"
)

var LastDFSVisited int

func DFS(target string) *TraceNode {
	LastDFSVisited = 0
	basic := map[string]struct{}{"Air": {}, "Earth": {}, "Fire": {}, "Water": {}, "Time": {}}
	cache := make(map[string]*TraceNode)

	// inline BFS fallback: cari satu trace via BFS jika DFS gagal
	bfsFallback := func(goal string) *TraceNode {
		// node struct for queue
		// type node struct{ trace *TraceNode }

		// init basic elements
		tqueue := []*TraceNode{}
		for b := range basic {
			tqueue = append(tqueue, &TraceNode{Product: b})
		}
		visited := map[string]bool{}
		for b := range basic {
			visited[b] = true
		}

		// memo for created nodes
		nodes := map[string]*TraceNode{}
		for _, n := range tqueue {
			nodes[n.Product] = n
		}

		// BFS
		for len(tqueue) > 0 {
			curr := tqueue[0]
			tqueue = tqueue[1:]
			prod := curr.Product

			if prod == goal {
				return curr
			}

			// explore all products that can be made
			for p, recs := range Graph {
				if visited[p] {
					continue
				}
				for _, r := range recs {
					a, b := r[0], r[1]
					if Tier[a] >= Tier[goal] && Tier[b] >= Tier[goal] {
						continue
					}
					if visited[a] && visited[b] {
						left, lok := nodes[a]
						right, rok := nodes[b]
						if lok && rok {
							newN := &TraceNode{Product: p, From: [2]string{a, b}, Parent: [2]*TraceNode{left, right}}
							visited[p] = true
							nodes[p] = newN
							tqueue = append(tqueue, newN)
							if p == goal {
								return newN
							}
							break
						}
					}
				}
			}
		}
		return nil
	}

	var dfsRec func(prod string, path map[string]bool) *TraceNode
	dfsRec = func(prod string, path map[string]bool) *TraceNode {
		LastDFSVisited++
		// cache sukses
		if n, ok := cache[prod]; ok {
			return n
		}
		// siklus deteksi
		if path[prod] {
			return nil
		}
		path[prod] = true
		defer delete(path, prod)

		// elemen dasar
		if _, ok := basic[prod]; ok {
			n := &TraceNode{Product: prod}
			cache[prod] = n
			return n
		}

		// coba setiap recipe
		for _, r := range Graph[prod] {
			a, b := r[0], r[1]
			if a == prod || b == prod {
				continue
			}
			// expand child a
			left := dfsRec(a, path)
			if left == nil {
				left = bfsFallback(a)
			}
			if left == nil || containsProduct(left, prod) {
				continue
			}

			// expand child b
			right := dfsRec(b, path)
			if right == nil {
				right = bfsFallback(b)
			}
			if right == nil || containsProduct(right, prod) {
				continue
			}

			n := &TraceNode{Product: prod, From: [2]string{a, b}, Parent: [2]*TraceNode{left, right}}
			cache[prod] = n
			return n
		}

		// gagal
		return nil
	}

	return dfsRec(target, make(map[string]bool))
}

// MultiDFS_Trace mengembalikan hingga maxResults jalur unik atau fallback ke MultiBFS
func MultiDFS_Trace(target string, maxResults int) []*TraceNode {
	// Reset the counter for this operation
	LastDFSVisited = 0

	var results []*TraceNode
	used := make(map[string]bool)
	var mu sync.Mutex

	recipeKey := func(a, b string) string {
		if a < b {
			return a + "+" + b
		}
		return b + "+" + a
	}

	// concurrency limiter and wait group
	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup

	for _, rec := range Graph[target] {
		a, b := rec[0], rec[1]
		// skip invalid or duplicate recipes
		if a == target || b == target || (Tier[a] >= Tier[target] && Tier[b] >= Tier[target]) {
			continue
		}

		mu.Lock()
		if len(results) >= maxResults {
			mu.Unlock()
			break
		}
		key := recipeKey(a, b)
		if used[key] {
			mu.Unlock()
			continue
		}
		used[key] = true
		mu.Unlock()

		// count this recipe attempt
		LastDFSVisited++

		// spawn goroutine to process this recipe
		sem <- struct{}{}
		wg.Add(1)
		go func(a, b string) {
			defer wg.Done()
			defer func() { <-sem }()

			left := DFS(a)
			right := DFS(b)
			if left == nil || right == nil {
				return
			}

			// exclude branches containing target to avoid cycles
			if containsProduct(left, target) || containsProduct(right, target) {
				return
			}

			// merge result
			n := &TraceNode{Product: target, From: [2]string{a, b}, Parent: [2]*TraceNode{left, right}}
			mu.Lock()
			if len(results) < maxResults {
				results = append(results, n)
			}
			mu.Unlock()
		}(a, b)
	}
	wg.Wait()
	return results
}

// helper: cek subtree memiliki produk di kedalaman apapun
func containsProduct(n *TraceNode, prod string) bool {
	if n == nil {
		return false
	}
	if n.Product == prod {
		return true
	}
	// cek seluruh parent
	for _, p := range n.Parent {
		if containsProduct(p, prod) {
			return true
		}
	}
	return false
}
