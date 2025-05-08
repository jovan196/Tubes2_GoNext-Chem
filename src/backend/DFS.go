package main

// DFS implements a depth-first search to build a single recipe trace for target
func DFS(target string) *TraceNode {
	// define basic elements set
	basic := map[string]struct{}{
		"Air": {}, "Water": {}, "Earth": {}, "Fire": {}, "Time": {},
	}

	// visited to avoid infinite recursion on cycles
	visited := make(map[string]bool)

	// recursive helper: returns TraceNode if product can be built
	var dfsRec func(string) *TraceNode
	dfsRec = func(prod string) *TraceNode {
		if visited[prod] {
			return nil
		}
		visited[prod] = true

		// base case: basic element
		if _, ok := basic[prod]; ok {
			return &TraceNode{Product: prod}
		}

		// try each two-input recipe for this product
		for _, rec := range Graph[prod] {
			left := dfsRec(rec[0])
			if left == nil {
				continue
			}
			right := dfsRec(rec[1])
			if right == nil {
				continue
			}

			// both ingredients built => create TraceNode
			return &TraceNode{
				Product: prod,
				From:    [2]string{rec[0], rec[1]},
				Parent:  [2]*TraceNode{left, right},
			}
		}

		// no recipe works
		return nil
	}

	return dfsRec(target)
}

func MultiDFS_Trace(target string, max int) []*TraceNode {
	var results []*TraceNode

	// For each direct recipe of target, attempt DFS for both ingredients
	for _, rec := range Graph[target] {
		left := DFS(rec[0])
		right := DFS(rec[1])
		if left != nil && right != nil {
			node := &TraceNode{
				Product: target,
				From:    [2]string{rec[0], rec[1]},
				Parent:  [2]*TraceNode{left, right},
			}
			results = append(results, node)
			if len(results) >= max {
				break
			}
		}
	}
	return results
}
