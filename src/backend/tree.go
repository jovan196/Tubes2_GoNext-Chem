package main

import "fmt"

type OutputNode struct {
	Name     string        `json:"name"`
	Children []*OutputNode `json:"children,omitempty"`
}

func toOutputTree(n *TraceNode) *OutputNode {
	if n == nil {
		return nil
	}
	if n.Parent[0] == nil && n.Parent[1] == nil {
		return &OutputNode{Name: n.Product}
	}
	return &OutputNode{
		Name: n.Product,
		Children: []*OutputNode{
			toOutputTree(n.Parent[0]),
			toOutputTree(n.Parent[1]),
		},
	}
}

func convertTraceToOutputRecursive(n *TraceNode) *OutputNode {
	if n == nil {
		return nil
	}
	node := &OutputNode{Name: n.Product}

	if n.Parent[0] != nil && n.Parent[1] != nil {
		node.Children = []*OutputNode{
			convertTraceToOutputRecursive(n.Parent[0]),
			convertTraceToOutputRecursive(n.Parent[1]),
		}
	}
	return node
}

func mergeTraceTrees(roots []*TraceNode) *OutputNode {
	if len(roots) == 0 {
		return nil
	}
	output := &OutputNode{Name: roots[0].Product}
	seen := make(map[string]bool)

	for i, root := range roots {
		hash := hashSubtree(root)
		if !seen[hash] {
			output.Children = append(output.Children, &OutputNode{
				Name:     fmt.Sprintf("#%d", i+1),
				Children: []*OutputNode{convertTraceToOutputRecursive(root)},
			})
			seen[hash] = true
		}
	}
	return output
}
