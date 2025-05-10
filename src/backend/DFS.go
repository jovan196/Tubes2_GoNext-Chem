package main

var LastDFSVisited int

/*func DFS(target string) *TraceNode {
	// basic elements
	basic := map[string]struct{}{"Air": {}, "Earth": {}, "Fire": {}, "Water": {}, "Time": {}}
	cache := make(map[string]*TraceNode)

	var dfsRec func(prod string, path map[string]bool) *TraceNode
	dfsRec = func(prod string, path map[string]bool) *TraceNode {
		// memoization hit (only success cached)
		if node, ok := cache[prod]; ok {
			return node
		}
		// cycle detection
		if path[prod] {
			return nil
		}
		path[prod] = true
		defer delete(path, prod)

		// base case: basic element
		if _, ok := basic[prod]; ok {
			n := &TraceNode{Product: prod}
			cache[prod] = n
			return n
		}

		// helper: DFS then BFS fallback for a child
		getChild := func(name string) *TraceNode {
			if child := dfsRec(name, path); child != nil {
				return child
			}
			return BFS(name)
		}

		// try each two-ingredient recipe
		for _, rec := range Graph[prod] {
			// same ingredient case: expand once
			if rec[0] == rec[1] {
				child := getChild(rec[0])
				if child == nil {
					continue
				}
				n := &TraceNode{Product: prod, From: [2]string{rec[0], rec[1]}, Parent: [2]*TraceNode{child, child}}
				cache[prod] = n
				return n
			}

			// distinct ingredients
			left := getChild(rec[0])
			if left == nil {
				continue
			}
			right := getChild(rec[1])
			if right == nil {
				continue
			}

			n := &TraceNode{Product: prod, From: [2]string{rec[0], rec[1]}, Parent: [2]*TraceNode{left, right}}
			cache[prod] = n
			return n
		}

		// no valid recipe found here
		return nil
	}

	// run DFS, fallback to BFS for entire target if needed
	if root := dfsRec(target, make(map[string]bool)); root != nil {
		return root
	}
	return BFS(target)
}*/

// TraceNode & Graph tetap sama ---------------------------------------------

// DFS builds *one* valid recipe tree for target or returns nil.
// - path: mendeteksi siklus di stack
// - cache: memoize hanya node sukses, gagal tidak dicache agar dicoba ulang
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
			// expand child a
			left := dfsRec(a, path)
			if left == nil {
				left = bfsFallback(a)
			}
			if left == nil {
				continue
			}

			// expand child b
			right := dfsRec(b, path)
			if right == nil {
				right = bfsFallback(b)
			}
			if right == nil {
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
		if len(results) >= maxResults {
			break
		}
		key := recipeKey(rec[0], rec[1])
		if used[key] {
			continue
		}

		left := DFS(rec[0])
		right := DFS(rec[1])
		if left == nil || right == nil {
			continue
		}

		used[key] = true
		n := &TraceNode{Product: target, From: [2]string{rec[0], rec[1]}, Parent: [2]*TraceNode{left, right}}
		results = append(results, n)
	}

	if len(results) == 0 {
		// fallback: gunakan BFS multi-trace
		return MultiBFS_Trace(target, maxResults)
	}
	return results
}
