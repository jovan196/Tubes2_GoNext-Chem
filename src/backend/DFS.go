package main

// TraceNode & Graph tetap sama ---------------------------------------------

// DFS builds *one* valid recipe tree for target or returns nil.
//  - pathVisited : mendeteksi siklus di stack saat ini
//  - cache       : memoization sukses/gagal sehingga sub-pohon tak dihitung ulang
func DFS(target string) *TraceNode {
	// ← cukup sering dipakai, jadikan map literal
	basic := map[string]struct{}{
		"Air": {}, "Earth": {}, "Fire": {}, "Water": {}, "Time": {},
	}

	// cache: product → *TraceNode (nil artinya “tidak mungkin dibuat”)
	cache := make(map[string]*TraceNode)

	var dfsRec func(prod string, path map[string]bool) *TraceNode
	dfsRec = func(prod string, path map[string]bool) *TraceNode {
		// 1. memoisation hit
		if n, ok := cache[prod]; ok {
			return n
		}

		// 2. deteksi siklus di jalur ini
		if path[prod] {
			return nil
		}
		path[prod] = true        // push
		defer delete(path, prod) // pop saat keluar

		// 3. elemen dasar
		if _, ok := basic[prod]; ok {
			node := &TraceNode{Product: prod}
			cache[prod] = node
			return node
		}

		// 4. coba setiap resep dua-bahan
		for _, rec := range Graph[prod] { // jika Graph[prod] == nil akan dilewati
			left := dfsRec(rec[0], path)
			if left == nil {
				continue
			}
			right := dfsRec(rec[1], path)
			if right == nil {
				continue
			}
			node := &TraceNode{
				Product: prod,
				From:    [2]string{rec[0], rec[1]},
				Parent:  [2]*TraceNode{left, right},
			}
			cache[prod] = node
			return node
		}

		// 5. gagal total
		cache[prod] = nil
		return nil
	}

	return dfsRec(target, make(map[string]bool))
}

// MultiDFS_Trace mengembalikan ≤ max jalur unik untuk target
func MultiDFS_Trace(target string, max int) []*TraceNode {
	var results []*TraceNode
	for _, rec := range Graph[target] {
		if len(results) >= max {
			break
		}
		left := DFS(rec[0])
		right := DFS(rec[1])
		if left == nil || right == nil {
			continue
		}
		results = append(results, &TraceNode{
			Product: target,
			From:    [2]string{rec[0], rec[1]},
			Parent:  [2]*TraceNode{left, right},
		})
	}
	return results
}
