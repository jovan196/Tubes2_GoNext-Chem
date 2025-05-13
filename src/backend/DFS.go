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
			// skip if child tier is not lower than parent
			if Tier[a] >= Tier[prod] || Tier[b] >= Tier[prod] {
				continue
			}
			// expand child a via DFS
			left := dfsRec(a, path)
			if left == nil {
				continue
			}
			if containsProduct(left, prod) {
				continue
			}

			// expand child b via DFS
			right := dfsRec(b, path)
			if right == nil {
				continue
			}
			if containsProduct(right, prod) {
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

	// helper DFS multi satu level
	for _, rec := range Graph[target] {
		a, b := rec[0], rec[1]
		// skip if rec mentions target itself
		if a == target || b == target {
			continue
		}
		// skip if child tier is not lower than parent
		if Tier[a] >= Tier[target] || Tier[b] >= Tier[target] {
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
