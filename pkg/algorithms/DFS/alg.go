package dfs

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/node"
	"fmt"
)

type StackInterface interface {
	Len() int
	Push(*node.Node)
	Pop() (*node.Node, error)
	Peek() (*node.Node, error)
}

type DeepSearch struct {
	// для удобства храним граф в структуре
	graph graph.Graph
	// и цель тоже
	target *node.Node
	// рабочая память алгоритма
	stack StackInterface
	// список запрещённых вершин
	forbiddenMap map[int]*node.Node
}

func NewDeepSearch(gr graph.Graph, stack StackInterface) DeepSearch {
	return DeepSearch{
		graph:        gr,
		target:       nil,
		stack:        stack,
		forbiddenMap: map[int]*node.Node{},
	}
}

func (ds *DeepSearch) FindTarget(initial_vertex, target *node.Node) ([]*node.Node, error) {
	ds.target = target
	ds.stack.Push(initial_vertex)
	// флаг, что решение найдено
	decisionFlag := false
	for !decisionFlag {
		if ds.stack.Len() == 0 {
			return []*node.Node{}, fmt.Errorf("decision was not found")
		}
		v, _ := ds.stack.Peek()
		decisionFlag = ds.findDescendants(v)
	}

	path := ds.unpackStack()
	return path, nil
}

func (ds *DeepSearch) findDescendants(source *node.Node) bool {
	old_n := ds.stack.Len()
	for i, e := range ds.graph.Edges {
		if e.Start == source && e.Label != enums.Closed {
			// если ребро ведёт в запрещённую вершину, помечаем его как закрытое
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
	// если ни одной вершины не удалось добавить в стек, значит,
	// у этой вершины нет больше потомков
	if old_n == ds.stack.Len() {
		v, _ := ds.stack.Pop()
		ds.forbiddenMap[v.Number] = v
	}
	return false
}

func (ds DeepSearch) unpackStack() []*node.Node {
	n := ds.stack.Len()
	path := make([]*node.Node, n)
	for i := 0; i < n; i++ {
		v, _ := ds.stack.Pop()
		path[n-i-1] = v
	}
	return path
}
