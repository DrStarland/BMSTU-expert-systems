package dfs

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/vertex"
	"fmt"
)

type StackInterface interface {
	Len() int
	Push(*vertex.Vertex)
	Pop() (*vertex.Vertex, error)
	Peek() (*vertex.Vertex, error)
}

type DeepSearch struct {
	// для удобства храним граф в структуре
	graph graph.Graph
	// и цель тоже
	target *vertex.Vertex
	// рабочая память алгоритма
	stack StackInterface
	// список запрещённых вершин
	forbiddenMap map[int]*vertex.Vertex
}

func NewDeepSearch(gr graph.Graph, stack StackInterface) DeepSearch {
	return DeepSearch{
		graph:        gr,
		target:       nil,
		stack:        stack,
		forbiddenMap: map[int]*vertex.Vertex{},
	}
}

func (ds *DeepSearch) FindTarget(initial_vertex, target *vertex.Vertex) ([]*vertex.Vertex, error) {
	ds.target = target

	// path[0] = initial_vertex
	ds.stack.Push(initial_vertex)
	decisionFlag := false
	for !decisionFlag {
		if ds.stack.Len() == 0 {
			return []*vertex.Vertex{}, fmt.Errorf("decision was not found")
		}

		v, _ := ds.stack.Peek()
		decisionFlag = ds.findDescendants(v)
	}

	path := ds.unpackStack()
	return path, nil
}

func (ds *DeepSearch) findDescendants(source *vertex.Vertex) bool {
	old_n := ds.stack.Len()
	for i, e := range ds.graph.Edges {
		if e.Start == source && e.Label != enums.Closed {
			if _, ok := ds.forbiddenMap[e.End.Number]; ok {
				ds.graph.Edges[i].Label = enums.Closed
				continue
			}
			// наилучшая ветвь
			ds.stack.Push(e.End)
			if e.End == ds.target {
				return true
			} else {
				return false
			}
		}
	}
	if old_n == ds.stack.Len() {
		v, _ := ds.stack.Pop()
		ds.forbiddenMap[v.Number] = v
	}
	return false
}

func (ds DeepSearch) unpackStack() []*vertex.Vertex {
	n := ds.stack.Len()
	path := make([]*vertex.Vertex, n)
	for i := 0; i < n; i++ {
		v, _ := ds.stack.Pop()
		path[n-i-1] = v
	}
	return path
}
