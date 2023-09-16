package main

import (
	bfs "expert_systems/pkg/algorithms/BFS"
	dfs "expert_systems/pkg/algorithms/DFS"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/queue"
	"expert_systems/pkg/models/stack"
	"expert_systems/pkg/models/vertex"
	"log"
)

func main() {

	gr, err := graph.NewGraph(
		0, 1,
		1, 4,
		0, 2,
		0, 3,
		2, 4,
		2, 5,
		3, 6,
		7, 9,
		9, 8,
		4, 8,
	)
	// 	1, 2,
	// 	2, 3,
	// 	3, 4,
	// 	3, 6,
	// 	5, 6,
	// 	6, 7,
	// 	6, 8,
	// )
	if err != nil {
		log.Panicln(err)
	}

	stck := stack.NewStack[*vertex.Vertex]()
	log.Println("Поиск в глубину:")
	algDFS := dfs.NewDeepSearch(gr, stck)
	source := gr.Vertexes[0]
	target := gr.Vertexes[8]
	path, err := algDFS.FindTarget(source, target)

	log.Println(path, err)

	for _, v := range path {
		log.Println(v.Number)
	}

	log.Println("Поиск в ширину:")
	qu := queue.NewQueue[*vertex.Vertex]()
	stck2 := stack.NewStack[*vertex.Vertex]()
	algBFS := bfs.NewWideSearch(gr, qu, stck2)
	source = gr.Vertexes[0]
	target = gr.Vertexes[8]
	path, err = algBFS.FindTarget(source, target)

	log.Println(path, err)
	for _, v := range path {
		log.Println(v.Number)
	}

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
