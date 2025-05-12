package main

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
	var results []*TraceNode
	used := make(map[string]bool)
	recipeKey := func(a, b string) string {
		if a < b {
			return a + "+" + b
		}
		return b + "+" + a
	}

	// helper DFS multi satu level
	for _, rec := range Graph[target] {
		a, b := rec[0], rec[1]
		// ═══ skip resep yang menyebutkan target sendiri ═══
		if a == target || b == target {
			continue
		}

		// skip resep jika tier elemen a dan b lebih tinggi dari target
		if Tier[a] >= Tier[target] && Tier[b] >= Tier[target] {
			continue
		}

		if len(results) >= maxResults {
			break
		}
		key := recipeKey(a, b)
		if used[key] {
			continue
		}

		left := DFS(rec[0])
		right := DFS(rec[1])
		if left == nil || right == nil {
			continue
		}

		used[key] = true
		if containsProduct(left, target) || containsProduct(right, target) {
			continue
		}
		n := &TraceNode{Product: target, From: [2]string{rec[0], rec[1]}, Parent: [2]*TraceNode{left, right}}
		results = append(results, n)
	}

	if len(results) == 0 {
		// fallback: gunakan BFS multi-trace
		return MultiBFS_Trace(target, maxResults)
	}
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
