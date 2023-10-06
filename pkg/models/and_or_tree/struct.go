package and_or_tree

import (
	"expert_systems/pkg/models/edge"
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/rule"
	"fmt"
)

type Tree struct {
	// основа графа -- список рёбер
	Rules []rule.Rule
	// карта всех вершин по номерам
	Nodes map[int]*node.Node
}

type Path struct {
	// само содержание пути
	Data []interface{}
}

// введите структуру графа как список попарных номеров вершин, являющихся началом и концом ребра
// например, NewTree(1, 2, 2, 3, 3, 1)
func appendAndGet(mv *map[int]*node.Node, number int) *node.Node {
	v, ok := (*mv)[number]
	if ok {
		return v
	}
	v = node.NewNode(number)
	(*mv)[number] = v
	return v
}

func NewTree(vertexNumbers ...int) (Tree, error) {
	n := len(vertexNumbers)
	if n%2 != 0 {
		return Tree{}, fmt.Errorf("incorrect graph: number of vertexex %d is odd", n)
	}

	// число рёбер в 2 раза меньше числа пар вершин
	gr := Tree{
		Edges: make([]edge.Edge, 0, n/2),
		Nodes: make(map[int]*node.Node),
	}
	for i := 0; i < n; i += 2 {
		gr.Edges = append(gr.Edges, edge.Edge{
			Start: appendAndGet(&gr.Nodes, vertexNumbers[i]),
			End:   appendAndGet(&gr.Nodes, vertexNumbers[i+1]),
			Label: enums.Open,
		})
	}

	return gr, nil
}
