package main

import (
	"fmt"
	"runtime"
	"sync"
)

var LastBFSVisited int

type TraceNode struct {
	Product string
	From    [2]string
	Parent  [2]*TraceNode
	Depth   int
}

// func convertTraceToOutput(n *TraceNode) *OutputNode {
// 	if n == nil {
// 		return nil
// 	}
// 	node := &OutputNode{Name: n.Product}

// 	if n.Parent[0] != nil && n.Parent[1] != nil {
// 		node.Children = []*OutputNode{
// 			convertTraceToOutput(n.Parent[0]),
// 			convertTraceToOutput(n.Parent[1]),
// 		}
// 	}

// 	return node
// }

func BFS(target string) *TraceNode {
	LastBFSVisited = 0
	buildableMemo = make(map[string]map[int]bool)
	if Tier[target] == 1 {
		return &TraceNode{Product: target}
	}

	visited := map[string]*TraceNode{}
	queue := []*TraceNode{}

	root := &TraceNode{Product: target}
	queue = append(queue, root)
	visited[target] = root

	for len(queue) > 0 {
		LastBFSVisited++
		curr := queue[0]
		queue = queue[1:]

		if Tier[curr.Product] == 1 {
			continue
		}

		for _, pair := range Graph[curr.Product] {

			a, b := pair[0], pair[1]
			if Tier[a] >= Tier[curr.Product] || Tier[b] >= Tier[curr.Product] {
				continue
			}

			// Check: pastikan a dan b bisa dibentuk dari basic
			if !canBuild(a, Tier[a]) || !canBuild(b, Tier[b]) {
				continue
			}
			left := visited[a]
			if left == nil {
				left = &TraceNode{Product: a}
				visited[a] = left
				queue = append(queue, left)
			}
			right := visited[b]
			if right == nil {
				right = &TraceNode{Product: b}
				visited[b] = right
				queue = append(queue, right)
			}

			curr.From = [2]string{a, b}
			curr.Parent = [2]*TraceNode{left, right}
			curr.Depth = 1 + max(left.Depth, right.Depth)
			break // hanya ambil 1 recipe valid
		}
	}

	// Validasi terakhir: jika root tidak punya Parent berarti tidak valid
	if root.Parent[0] == nil && root.Parent[1] == nil {
		if Tier[root.Product] != 1 { // bukan basic
			return nil
		}
	}
	return root
}

var buildableMemo = make(map[string]map[int]bool)
var buildableMutex sync.RWMutex

func canBuild(target string, tierLimit int) bool {
	return canBuildInternal(target, tierLimit, map[string]bool{})
}

func canBuildInternal(target string, tierLimit int, visited map[string]bool) bool {
	buildableMutex.RLock()
	if m, ok := buildableMemo[target]; ok {
		if val, ok := m[tierLimit]; ok {
			buildableMutex.RUnlock()
			return val
		}
	}
	buildableMutex.RUnlock()

	// Jika sedang diproses (recursive cycle)
	if visited[target] {
		return false
	}
	visited[target] = true

	if Tier[target] == 1 {
		buildableMutex.Lock()
		if _, ok := buildableMemo[target]; !ok {
			buildableMemo[target] = make(map[int]bool)
		}
		buildableMemo[target][tierLimit] = true
		buildableMutex.Unlock()
		return true
	}

	for _, pair := range Graph[target] {
		a, b := pair[0], pair[1]
		if Tier[a] >= tierLimit || Tier[b] >= tierLimit {
			continue
		}
		if canBuildInternal(a, tierLimit, copyMap(visited)) && canBuildInternal(b, tierLimit, copyMap(visited)) {
			buildableMutex.Lock()
			if _, ok := buildableMemo[target]; !ok {
				buildableMemo[target] = make(map[int]bool)
			}
			buildableMemo[target][tierLimit] = true
			buildableMutex.Unlock()
			return true
		}
	}

	buildableMutex.Lock()
	if _, ok := buildableMemo[target]; !ok {
		buildableMemo[target] = make(map[int]bool)
	}
	buildableMemo[target][tierLimit] = false
	buildableMutex.Unlock()
	return false
}

func copyMap(original map[string]bool) map[string]bool {
	copy := make(map[string]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// func exists(m map[string]struct{}, key string) bool {
// 	_, ok := m[key]
// 	return ok
// }

func MultiBFS_Trace(target string, maxResults int) *OutputNode {
	LastBFSVisited = 0
	buildableMemo = make(map[string]map[int]bool)
	if Tier[target] == 1 {
		return &OutputNode{Name: target}
	}

	// type Work struct {
	// 	Node *TraceNodeMulti
	// }

	var (
		queue      = []*TraceNodeMulti{{Product: target}}
		roots      []*TraceNodeMulti
		seenHashes = make(map[string]bool)
		mu         sync.Mutex
		wg         sync.WaitGroup
		sem        = make(chan struct{}, runtime.NumCPU())
	)

	for len(queue) > 0 && len(roots) < maxResults {
		curr := queue[0]
		queue = queue[1:]

		if Tier[curr.Product] == 1 {
			continue
		}

		for _, pair := range Graph[curr.Product] {
			a, b := pair[0], pair[1]
			LastBFSVisited++
			if Tier[a] >= Tier[curr.Product] || Tier[b] >= Tier[curr.Product] {
				continue
			}
			if !canBuild(a, Tier[a]) || !canBuild(b, Tier[b]) {
				continue
			}

			sem <- struct{}{}
			wg.Add(1)
			go func(a, b string, curr *TraceNodeMulti) {
				defer wg.Done()
				defer func() { <-sem }()

				left := &TraceNodeMulti{Product: a}
				right := &TraceNodeMulti{Product: b}
				node := &TraceNodeMulti{
					Product: curr.Product,
					From:    [2]string{a, b},
					Parent:  [2]*TraceNodeMulti{left, right},
					Depth:   1 + max(left.Depth, right.Depth),
				}

				expand(left)
				expand(right)

				h := hashSubtreeMulti(node)

				mu.Lock()
				if !seenHashes[h] && len(roots) < maxResults {
					seenHashes[h] = true
					roots = append(roots, node)
					queue = append(queue, node)
				}
				mu.Unlock()
			}(a, b, curr)
		}

		wg.Wait()
		if len(roots) == 0 {
			return nil
		}
	}

	return mergeTraceTreesMulti(roots)
}

func expand(n *TraceNodeMulti) {
	if n == nil || (n.Parent[0] != nil && n.Parent[1] != nil) || Tier[n.Product] == 1 {
		return
	}
	for _, pair := range Graph[n.Product] {
		a, b := pair[0], pair[1]
		if Tier[a] >= Tier[n.Product] || Tier[b] >= Tier[n.Product] {
			continue
		}
		if !canBuild(a, Tier[a]) || !canBuild(b, Tier[b]) {
			continue
		}
		left := &TraceNodeMulti{Product: a}
		right := &TraceNodeMulti{Product: b}
		n.From = [2]string{a, b}
		n.Parent = [2]*TraceNodeMulti{left, right}
		n.Depth = 1 + max(left.Depth, right.Depth)

		expand(left)
		expand(right)
		break
	}
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

func hashSubtreeMulti(n *TraceNodeMulti) string {
	if n == nil {
		return ""
	}
	if n.Parent[0] == nil && n.Parent[1] == nil {
		return n.Product
	}
	l := hashSubtreeMulti(n.Parent[0])
	r := hashSubtreeMulti(n.Parent[1])
	if l > r {
		l, r = r, l
	}
	return fmt.Sprintf("%s(%s+%s)", n.Product, l, r)
}

// func containProduct(n *TraceNode, prod string) bool {
// 	if n == nil {
// 		return false
// 	}
// 	if n.Product == prod {
// 		return true
// 	}
// 	return containProduct(n.Parent[0], prod) ||
// 		containProduct(n.Parent[1], prod)
// }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
