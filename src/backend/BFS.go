package main

func BFS(target string) ([][2]string, bool) {
	type Node struct {
		Element string
		Path    [][2]string
	}

	queue := []Node{
		{Element: "Air", Path: nil},
		{Element: "Water", Path: nil},
		{Element: "Earth", Path: nil},
		{Element: "Fire", Path: nil},
	}

	visited := map[string]bool{
		"Air": true, "Water": true, "Earth": true, "Fire": true,
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for product, recipes := range Graph {
			for _, pair := range recipes {
				if (curr.Element == pair[0] && visited[pair[1]]) || (curr.Element == pair[1] && visited[pair[0]]) {
					if product == target {
						return append(curr.Path, pair), true
					}
					if !visited[product] {
						visited[product] = true
						queue = append(queue, Node{
							Element: product,
							Path:    append(curr.Path, pair),
						})
					}
				}
			}
		}
	}

	return nil, false
}
