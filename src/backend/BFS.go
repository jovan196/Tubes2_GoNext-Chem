package main

import "sync"

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

func MultiBFS(target string, max int) [][][2]string {
	type Node struct {
		Element string
		Path    [][2]string
	}

	type Result struct {
		Path [][2]string
	}

	var results [][][2]string
	resultChan := make(chan [][2]string, max)
	var wg sync.WaitGroup
	var mu sync.Mutex

	visitedGlobal := map[string]bool{
		"Air": true, "Water": true, "Earth": true, "Fire": true,
	}

	for product, recipes := range Graph {
		if product == target {
			for _, pair := range recipes {
				wg.Add(1)
				go func(pair [2]string) {
					defer wg.Done()

					queue := []Node{
						{Element: pair[0], Path: [][2]string{}},
						{Element: pair[1], Path: [][2]string{}},
					}

					visited := map[string]bool{}
					for k, v := range visitedGlobal {
						visited[k] = v
					}

					visited[pair[0]] = true
					visited[pair[1]] = true

					for len(queue) > 0 {
						curr := queue[0]
						queue = queue[1:]

						for prod, recs := range Graph {
							for _, rec := range recs {
								if (curr.Element == rec[0] && visited[rec[1]]) || (curr.Element == rec[1] && visited[rec[0]]) {
									if prod == target {
										result := append(curr.Path, rec)
										mu.Lock()
										if len(results) < max {
											results = append(results, result)
											resultChan <- result
										}
										mu.Unlock()
										return
									}
									if !visited[prod] {
										visited[prod] = true
										queue = append(queue, Node{
											Element: prod,
											Path:    append(curr.Path, rec),
										})
									}
								}
							}
						}
					}
				}(pair)
			}
		}
	}

	// Tunggu semua goroutine selesai
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Kumpulkan hasil dari channel
	finalResults := [][][2]string{}
	for path := range resultChan {
		finalResults = append(finalResults, path)
		if len(finalResults) >= max {
			break
		}
	}

	return finalResults
}
