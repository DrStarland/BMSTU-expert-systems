package bfs

import (
	dfs "expert_systems/pkg/algorithms/DFS"
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/vertex"
	"fmt"
)

type QueueInterface interface {
	Len() int
	Push(*vertex.Vertex)
	Pop() (*vertex.Vertex, error)
	Peek() (*vertex.Vertex, error)
}

type WideSearch struct {
	// для удобства храним граф в структуре
	Graph graph.Graph
	// и цель тоже
	target *vertex.Vertex
	// рабочая память алгоритма
	queue     QueueInterface
	pathStack dfs.StackInterface
	// список запрещённых вершин
	forbiddenMap map[int]*vertex.Vertex
}

func NewWideSearch(gr graph.Graph, queue QueueInterface, stack dfs.StackInterface) WideSearch {
	return WideSearch{
		Graph:        gr,
		target:       nil,
		queue:        queue,
		pathStack:    stack,
		forbiddenMap: map[int]*vertex.Vertex{},
	}
}

func (bs *WideSearch) FindTarget(initial_vertex, target *vertex.Vertex) ([]*vertex.Vertex, error) {
	bs.target = target
	path_hidden := make([]*vertex.Vertex, 0)

	// path[0] = initial_vertex
	bs.queue.Push(initial_vertex)
	decisionFlag := false
	for !decisionFlag {
		if bs.queue.Len() == 0 {
			return []*vertex.Vertex{}, fmt.Errorf("decision was not found")
		}

		v, _ := bs.queue.Pop()
		decisionFlag = bs.findDescendants(v, &path_hidden)
	}

	path2 := bs.unpackQueue()
	return path2, nil
}

func (bs *WideSearch) findDescendants(source *vertex.Vertex, path *[]*vertex.Vertex) bool {
	old_n := bs.queue.Len()
	for i, e := range bs.Graph.Edges {
		if e.Start == source && e.Label != enums.Closed {
			if _, ok := bs.forbiddenMap[e.End.Number]; ok {
				bs.Graph.Edges[i].Label = enums.Closed
				continue
			}
			// наилучшая ветвь
			bs.queue.Push(e.End)
			if e.End == bs.target {
				return true
			}
		}
	}
	if old_n == bs.queue.Len() {
		bs.forbiddenMap[source.Number] = source
	}
	return false
}

func (bs WideSearch) unpackQueue() []*vertex.Vertex {
	n := bs.queue.Len()
	path := make([]*vertex.Vertex, n)
	for i := 0; i < n; i++ {
		v, _ := bs.queue.Pop()
		path[n-i-1] = v
	}
	return path
}
