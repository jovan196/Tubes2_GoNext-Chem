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

func convertTraceToOutput(n *TraceNode) *OutputNode {
	if n == nil {
		return nil
	}
	node := &OutputNode{Name: n.Product}

	if n.Parent[0] != nil && n.Parent[1] != nil {
		node.Children = []*OutputNode{
			convertTraceToOutput(n.Parent[0]),
			convertTraceToOutput(n.Parent[1]),
		}
	}

	return node
}

func BFS(target string) *TraceNode {
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
		curr := queue[0]
		fmt.Println(curr.Product)
		queue = queue[1:]

		if Tier[curr.Product] == 1 {
			continue
		}

		for _, pair := range Graph[curr.Product] {

			a, b := pair[0], pair[1]
			fmt.Printf("%s %s\n", a, b)
			if Tier[a] >= Tier[curr.Product] || Tier[b] >= Tier[curr.Product] {
				fmt.Println("giganig")
				continue
			}
			fmt.Println(1)

			// ⛔ Check: pastikan a dan b bisa dibentuk dari basic
			if !canBuild(a, Tier[a]) || !canBuild(b, Tier[b]) {
				continue
			}
			fmt.Println(2)

			left := visited[a]
			if left == nil {
				left = &TraceNode{Product: a}
				visited[a] = left
				queue = append(queue, left)
			}
			fmt.Println(3)
			right := visited[b]
			if right == nil {
				right = &TraceNode{Product: b}
				visited[b] = right
				queue = append(queue, right)
			}
			fmt.Println(4)

			curr.From = [2]string{a, b}
			curr.Parent = [2]*TraceNode{left, right}
			curr.Depth = 1 + max(left.Depth, right.Depth)
			fmt.Println(5)
			break // hanya ambil 1 recipe valid
		}
	}

	// Validasi terakhir: jika root tidak punya Parent berarti tidak valid
	if root.Parent[0] == nil && root.Parent[1] == nil {
		if Tier[root.Product] != 1 { // bukan basic
			return nil
		}
	}
	fmt.Println("basing")
	return root
}

var buildableMemo = make(map[string]map[int]bool)
var buildableMutex sync.Mutex

func canBuild(target string, tierLimit int) bool {
	buildableMutex.Lock()
	if m, ok := buildableMemo[target]; ok {
		if val, ok := m[tierLimit]; ok {
			buildableMutex.Unlock()
			return val
		}
	} else {
		buildableMemo[target] = make(map[int]bool)
	}
	buildableMutex.Unlock()

	if Tier[target] == 1 {
		buildableMutex.Lock()
		buildableMemo[target][tierLimit] = true
		buildableMutex.Unlock()
		return true
	}

	for _, pair := range Graph[target] {
		a, b := pair[0], pair[1]
		if Tier[a] >= tierLimit || Tier[b] >= tierLimit {
			continue
		}
		if canBuild(a, tierLimit) && canBuild(b, tierLimit) {
			buildableMutex.Lock()
			buildableMemo[target][tierLimit] = true
			buildableMutex.Unlock()
			return true
		}
	}

	buildableMutex.Lock()
	buildableMemo[target][tierLimit] = false
	buildableMutex.Unlock()
	return false
}

func exists(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}

func MultiBFS_Trace(target string, maxResults int) *OutputNode {
	LastBFSVisited = 0
	basic := []string{"Air", "Water", "Earth", "Fire", "Time"}
	queue := [][]*TraceNode{}
	nodes := make(map[string][]*TraceNode)
	seenHash := make(map[string]bool)
	var counter int32
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())

	for _, b := range basic {
		n := &TraceNode{Product: b}
		nodes[b] = []*TraceNode{n}
		queue = append(queue, []*TraceNode{n})
	}

	roots := []*TraceNode{}

	for len(queue) > 0 && int(atomic.LoadInt32(&counter)) < maxResults {
		currLevel := queue[0]
		queue = queue[1:]
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
							go func(ta, tb *TraceNode, out string) {
								defer wg.Done()
								defer func() { <-sem }()

								if containProduct(ta, out) || containProduct(tb, out) {
									return
								}

								n := &TraceNode{
									Product: out,
									From:    [2]string{ta.Product, tb.Product},
									Parent:  [2]*TraceNode{ta, tb},
									Depth:   1 + max(ta.Depth, tb.Depth),
								}
								h := hashSubtree(n)

								mu.Lock()
								if seenHash[h] {
									mu.Unlock()
									return
								}
								seenHash[h] = true

								nodes[out] = append(nodes[out], n)
								if out == target && int(atomic.LoadInt32(&counter)) < maxResults {
									roots = append(roots, n)
									atomic.AddInt32(&counter, 1)
								}
								mu.Unlock()

								nextLevel = append(nextLevel, n)
							}(ta, tb, outProd)
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

	return mergeTraceTrees(roots)
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
