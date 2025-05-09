package main

// ----- util ---------------------------------------------------------------
func isBasic(e string) bool {
	switch e {
	case "Air", "Earth", "Fire", "Water", "Time":
		return true
	default:
		return false
	}
}

// ----- 1  ×  DFS (satu jalur) --------------------------------------------
func DFS(target string) *TraceNode {
	if isBasic(target) {
		return &TraceNode{Product: target}
	}

	// memo      : hasil yg sudah dihitung (termasuk `nil` utk dead-end)
	// onStack   : deteksi siklus pada jalur saat ini
	memo := map[string]*TraceNode{}
	onStack := map[string]bool{}

	var rec func(string) *TraceNode
	rec = func(prod string) *TraceNode {
		if node, ok := memo[prod]; ok { // sudah dihitung sebelumnya
			return node
		}
		if isBasic(prod) {
			n := &TraceNode{Product: prod}
			memo[prod] = n
			return n
		}
		if onStack[prod] { // siklus → dead-end
			return nil
		}
		onStack[prod] = true
		defer delete(onStack, prod)

		for _, pair := range Graph[prod] { // coba setiap resep
			l := rec(pair[0])
			if l == nil {
				continue
			}
			r := rec(pair[1])
			if r == nil {
				continue
			}
			n := &TraceNode{
				Product: prod,
				From:    [2]string{pair[0], pair[1]},
				Parent:  [2]*TraceNode{l, r},
			}
			memo[prod] = n
			return n
		}
		memo[prod] = nil // dead-end disimpan supaya tidak diulang
		return nil
	}

	return rec(target)
}

// ----- Multi-DFS : k jalur berbeda ---------------------------------------
func MultiDFS_Trace(target string, max int) []*TraceNode {
	// jika elemen dasar, kembalikan node tunggal
	if isBasic(target) {
		return []*TraceNode{{Product: target}}
	}

	var res []*TraceNode
	for _, pair := range Graph[target] {
		left := DFS(pair[0])
		right := DFS(pair[1])
		if left != nil && right != nil {
			res = append(res, &TraceNode{
				Product: target,
				From:    [2]string{pair[0], pair[1]},
				Parent:  [2]*TraceNode{left, right},
			})
			if len(res) >= max {
				break
			}
		}
	}
	return res
}
