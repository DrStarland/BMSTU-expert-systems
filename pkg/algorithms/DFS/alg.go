package dfs

import (
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/vertex"
)

// Using Generics to define Type in Stake to Use Structs, too.
type StackInterface interface {
	Push(vertex.Vertex)
	Pop() (vertex.Vertex, error)
	Peek() (vertex.Vertex, error)
}

type DeepSearch struct {
	// для удобства храним граф в структуре
	graph graph.Graph
	// и цель тоже
	target vertex.Vertex
	// рабочая память алгоритма
	st StackInterface
	// список запрещённых вершин
	forbiddenList []*vertex.Vertex
}

func NewDeepSearch(gr graph.Graph, stack StackInterface) DeepSearch {
	return DeepSearch{
		graph:         gr,
		target:        vertex.Vertex{},
		st:            stack,
		forbiddenList: []*vertex.Vertex{},
	}
}

func (ds DeepSearch) FindTarget(initial_vertex, target vertex.Vertex) {

}

func (ds DeepSearch) FindDescendants(vr *vertex.Vertex) {
	for _, v := range ds.graph {
		if v.Start == vr {
			return
		}
	}
}
