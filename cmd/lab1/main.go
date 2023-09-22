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
	if err != nil {
		log.Panicln(err)
	}

	stck := stack.NewStack[*vertex.Vertex]()
	log.Println("Поиск в глубину:")
	algDFS := dfs.NewDeepSearch(gr, stck)
	source := gr.Vertexes[1]
	target := gr.Vertexes[8]
	path, err := algDFS.FindTarget(source, target)

	log.Println(path, err)

	for _, v := range path {
		log.Println(v.Number)
	}

	log.Println("Поиск в ширину:")
	qu := queue.NewQueue[*vertex.Vertex]()
	algBFS := bfs.NewWideSearch(gr, qu)
	source = gr.Vertexes[1]
	target = gr.Vertexes[8]
	path, err = algBFS.FindTarget(source, target)

	log.Println(path, err)
	for _, v := range path {
		log.Println(v.Number)
	}
}
