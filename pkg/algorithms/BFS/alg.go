package bfs

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/graph"
	"expert_systems/pkg/models/node"
	"fmt"
)

type QueueInterface interface {
	Len() int
	Push(*node.Node)
	Pop() (*node.Node, error)
	Peek() (*node.Node, error)
}

type WideSearch struct {
	// для удобства храним граф в структуре
	Graph graph.Graph
	// и цель тоже
	target *node.Node
	// рабочая память алгоритма
	queue QueueInterface
	// вспомогательная переменная для хранения путей к каждой вершине
	path_hidden map[int][]*node.Node
	// список запрещённых вершин
	forbiddenMap map[int]*node.Node
}

func NewWideSearch(gr graph.Graph, queue QueueInterface) WideSearch {
	return WideSearch{
		Graph:        gr,
		target:       nil,
		queue:        queue,
		path_hidden:  map[int][]*node.Node{},
		forbiddenMap: map[int]*node.Node{},
	}
}

func (bs *WideSearch) FindTarget(initial_node, target *node.Node) ([]*node.Node, error) {
	bs.target = target

	bs.queue.Push(initial_node)
	// флаг, что решение найдено
	decisionFlag := false
	for !decisionFlag {
		if bs.queue.Len() == 0 {
			return []*node.Node{}, fmt.Errorf("decision was not found")
		}

		v, _ := bs.queue.Pop()
		decisionFlag = bs.findDescendants(v)
	}

	return bs.path_hidden[bs.target.Number], nil
}

func (bs *WideSearch) findDescendants(source *node.Node) bool {
	old_n := bs.queue.Len()
	for i, e := range bs.Graph.Edges {
		if e.Start == source && e.Label != enums.Closed {
			// если ребро ведёт в запрещённую вершину, помечаем его как закрытое
			if _, ok := bs.forbiddenMap[e.End.Number]; ok {
				bs.Graph.Edges[i].Label = enums.Closed
				continue
			}

			if _, ok := bs.path_hidden[e.End.Number]; !ok {
				// копируем путь до текущей вершины и добавляем текущую -- это будет
				// путём до следующей вершины этого ребра
				temp := make([]*node.Node, len(bs.path_hidden[e.Start.Number]))
				copy(temp, bs.path_hidden[e.Start.Number])
				bs.path_hidden[e.End.Number] = append(temp, e.Start)
			}

			bs.queue.Push(e.End)
			if e.End == bs.target {
				bs.path_hidden[e.End.Number] = append(bs.path_hidden[e.End.Number], e.End)
				return true
			}
		}
	}
	// если ни одной вершины не удалось добавить в очередь, значит,
	// у этой вершины нет больше потомков
	if old_n == bs.queue.Len() {
		bs.forbiddenMap[source.Number] = source
	}
	return false
}
