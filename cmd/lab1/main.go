package main

import (
	dfs "expert_systems/pkg/algorithms/DFS"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/stack"
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

	stck := stack.NewStack[*node.Node]()
	log.Println("Поиск в глубину:")
	algDFS := dfs.NewDeepSearch(gr, stck)
	source := gr.Nodes[0]
	target := gr.Nodes[8]
	path, err := algDFS.FindTarget(source, target)

	log.Println(path, err)

	for _, v := range path {
		log.Println(v.Number)
	}
}
