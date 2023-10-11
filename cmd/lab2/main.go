package main

import (
	bfs "expert_systems/pkg/algorithms/BFS"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/queue"
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

	log.Println("Поиск в ширину:")
	qu := queue.NewQueue[*node.Node]()
	algBFS := bfs.NewWideSearch(gr, qu)
	source := gr.Nodes[1]
	target := gr.Nodes[8]
	path, err := algBFS.FindTarget(source, target)

	log.Println(path, err)
	for _, v := range path {
		log.Println(v.Number)
	}
}
