package graph

import (
	"expert_systems/pkg/models/edge"
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/node"
	"fmt"
)

type Graph struct {
	// основа графа -- список рёбер
	Edges []edge.Edge
	// карта всех вершин по номерам
	Nodes map[int]*node.Node
}

type Path struct {
	// само содержание пути
	Data []*node.Node
}

// введите структуру графа как список попарных номеров вершин, являющихся началом и концом ребра
// например, NewGraph(1, 2, 2, 3, 3, 1)
func appendAndGet(mv *map[int]*node.Node, number int) *node.Node {
	v, ok := (*mv)[number]
	if ok {
		return v
	}
	v = node.NewNode(number)
	(*mv)[number] = v
	return v
}

func NewGraph(nodeNumbers ...int) (Graph, error) {
	n := len(nodeNumbers)
	if n%2 != 0 {
		return Graph{}, fmt.Errorf("incorrect graph: number of nodeex %d is odd", n)
	}

	// число рёбер в 2 раза меньше числа пар вершин
	gr := Graph{
		Edges: make([]edge.Edge, 0, n/2),
		Nodes: make(map[int]*node.Node),
	}
	for i := 0; i < n; i += 2 {
		gr.Edges = append(gr.Edges, edge.Edge{
			Start: appendAndGet(&gr.Nodes, nodeNumbers[i]),
			End:   appendAndGet(&gr.Nodes, nodeNumbers[i+1]),
			Label: enums.Open,
		})
	}

	return gr, nil
}
