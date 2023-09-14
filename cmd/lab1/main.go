package main

import (
	"expert_systems/pkg/models/graph"
)

func main() {

	gr, err := graph.NewGraph(
		1, 2,
		2, 3,
		3, 4,
		3, 6,
		5, 6,
		6, 7,
		6, 8,
	)

	// alg = AlgorithmDFS(graph)
	// source = Vertex(1)
	// target = Vertex(8)

	// path = alg.search(source, target)

	// utils.print_path(path)
	// utils.show_graph(graph, path, node_color_map={
	// 	source.number: 'red',
	// 	target.number: 'red',
	// })

	// stck := stack.NewStack[string]()
	// stck.Push("1")
	// stck2 := stack.NewStack[[]int]()
	// stck2.Push([]int{1, 2, 3, 4, 5})
	// element, _ := stck2.Peek()

	// log.Println(element)
}
