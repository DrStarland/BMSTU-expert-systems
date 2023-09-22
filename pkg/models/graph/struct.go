package graph

import (
	"expert_systems/pkg/models/edge"
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/vertex"
	"fmt"
)

type Graph struct {
	// основа графа -- список рёбер
	Edges []edge.Edge
	// карта всех вершин по номерам
	Vertexes map[int]*vertex.Vertex
}

// введите структуру графа как список попарных номеров вершин, являющихся началом и концом ребра
// например, NewGraph(1, 2, 2, 3, 3, 1)
func appendAndGet(mv *map[int]*vertex.Vertex, number int) *vertex.Vertex {
	v, ok := (*mv)[number]
	if ok {
		return v
	}
	v = vertex.NewVertex(number)
	(*mv)[number] = v
	return v
}

func NewGraph(vertexNumbers ...int) (Graph, error) {
	n := len(vertexNumbers)
	if n%2 != 0 {
		return Graph{}, fmt.Errorf("incorrect graph: number of vertexex %d is odd", n)
	}

	// число рёбер в 2 раза меньше числа пар вершин
	gr := Graph{
		Edges:    make([]edge.Edge, 0, n/2),
		Vertexes: make(map[int]*vertex.Vertex),
	}
	for i := 0; i < n; i += 2 {
		gr.Edges = append(gr.Edges, edge.Edge{
			Start: appendAndGet(&gr.Vertexes, vertexNumbers[i]),
			End:   appendAndGet(&gr.Vertexes, vertexNumbers[i+1]),
			Label: enums.Open,
		})
	}

	return gr, nil
}
