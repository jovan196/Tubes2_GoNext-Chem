package main

// elemen dasar
var basic = []string{"Air", "Earth", "Fire", "Water"}

type TraceNode struct {
	Product string
	From    [2]string
	Parent  [2]*TraceNode
	Depth   int
}

/*func isBasic(e string) bool {
	for _, b := range basic {
		if b == e {
			return true
		}
	}
	return false
}*/

// ----------  SINGLE  BFS  ----------
func BFS(target string) *TraceNode {
	all := map[string]*TraceNode{}
	queue := []*TraceNode{}
	visited := map[string]bool{}

	// init dasar
	for _, b := range basic {
		n := &TraceNode{Product: b, Depth: 0}
		all[b] = n
		queue = append(queue, n)
		if b == target {
			return n // <--  langsung selesai
		}

	}

	for len(queue) > 0 {
		curr := queue[0]  // <-- ambil dulu
		queue = queue[1:] // pop
		visited[curr.Product] = true

		for prod, recs := range Graph {
			if visited[prod] { // boleh skip hanya jika node SUDAH diproses
				continue
			}
			for _, pair := range recs {
				a, b := pair[0], pair[1]
				left, okL := all[a]
				right, okR := all[b]
				if !okL || !okR { // kedua bahan belum tersedia
					continue
				}

				node := &TraceNode{
					Product: prod,
					From:    [2]string{a, b},
					Parent:  [2]*TraceNode{left, right},
					Depth:   max(left.Depth, right.Depth) + 1,
				}
				all[prod] = node
				if prod == target {
					return node
				}
				queue = append(queue, node)
			}
		}
	}
	return nil
}

// ----------  MULTIPLE  BFS  ----------
func MultiBFS_Trace(target string, maxSol int) []*TraceNode {
	var solutions []*TraceNode
	for _, pair := range Graph[target] { // coba tiap resep langsung
		startA, startB := pair[0], pair[1]

		// fresh structures utk tiap percobaan
		all := map[string]*TraceNode{
			startA: {Product: startA},
			startB: {Product: startB},
		}
		queue := []*TraceNode{all[startA], all[startB]}
		visited := map[string]bool{}
		for _, b := range basic {
			visited[b] = true
		}

		for len(queue) > 0 && len(solutions) < maxSol {
			curr := queue[0]
			queue = queue[1:]
			visited[curr.Product] = true

			for prod, recs := range Graph {
				if visited[prod] {
					continue
				}
				for _, rc := range recs {
					a, b := rc[0], rc[1]
					left, okL := all[a]
					right, okR := all[b]
					if !okL || !okR {
						continue
					}
					node := &TraceNode{
						Product: prod,
						From:    [2]string{a, b},
						Parent:  [2]*TraceNode{left, right},
						Depth:   max(left.Depth, right.Depth) + 1,
					}
					all[prod] = node
					if prod == target {
						solutions = append(solutions, node)
						break // temukan satu jalur; cari jalur lain dg resep berbeda
					}
					queue = append(queue, node)
				}
			}
		}
		if len(solutions) >= maxSol {
			break
		}
	}
	return solutions
}

// ----------  helper ----------
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
