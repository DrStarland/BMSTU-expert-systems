package dfs_tree

import (
	"expert_systems/pkg/models/and_or_tree"
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/rule"
	"expert_systems/pkg/models/stack"
	"log"
)

type StackInterface interface {
	Len() int
	Push(rule.Rule)
	Pop() (rule.Rule, error)
	Peek() (rule.Rule, error)
}

type DeepSearch struct {
	// для удобства храним граф в структуре
	knowledgebase and_or_tree.Tree
	// и цель тоже
	target *node.Node
	// рабочая память алгоритма
	stack StackInterface

	openNodes        []*node.Node
	openRules        []rule.Rule
	closedNodes      map[int]*node.Node
	closedNodesOrder []int
	closedRules      map[int]rule.Rule
	closedRulesOrder []int
	path             StackInterface

	// // список запрещённых вершин
	// forbiddenNodesMap map[int]*node.Node
}

func NewSearch(tr and_or_tree.Tree) DeepSearch {
	stck := stack.NewStack[rule.Rule]()
	return DeepSearch{
		knowledgebase: tr,
		target:        nil,
		stack:         stck,
		openNodes:     []*node.Node{},
		openRules:     []rule.Rule{},
		closedNodes:   map[int]*node.Node{},
		closedRules:   map[int]rule.Rule{},
		path:          nil,
	}
}

func (ds *DeepSearch) init(initial_nodes []*node.Node, target *node.Node) error {
	ds.target = target

	for _, nod := range initial_nodes {
		ds.closedNodes[nod.Number] = nod
		ds.closedNodesOrder = append(ds.closedNodesOrder, nod.Number)
	}

	for k, nod := range ds.knowledgebase.Nodes {
		if _, ok := ds.closedNodes[k]; !ok {
			ds.openNodes = append(ds.openNodes, nod)
		}
	}

	log.Println("Knowledgebase rules: ", ds.knowledgebase.Rules)
	ds.openRules = append(ds.openRules, ds.knowledgebase.Rules...)

	for ds.stack.Len() != 0 {
		ds.stack.Pop()
	}
	return nil
}

func (ds *DeepSearch) FindTarget(target *node.Node, inputs ...*node.Node) ([]*node.Node, error) {
	ds.init(inputs, target)
	log.Println("Target: ", target.Number)
	// флаг, что решение найдено
	decisionFlag := false
	decisionCanBeFound := true
	log.Println(ds.openRules)
	for !decisionFlag && decisionCanBeFound {
		// if ds.stack.Len() == 0 {
		// 	return []*node.Node{}, fmt.Errorf("decision was not found")
		// }
		// _, _ := ds.stack.Peek()
		decisionFlag, decisionCanBeFound = ds.findRules()
	}

	if !decisionCanBeFound {
		log.Println("Нет решения")
	}

	log.Println(ds.closedNodesOrder, ds.closedRulesOrder)
	// path := ds.unpackStack()
	return nil, nil
}

func (ds *DeepSearch) findRules() (bool, bool) {
	for i, r := range ds.openRules {
		flag := len(r.Inputs)
		for _, nod := range r.Inputs {
			if _, ok := ds.closedNodes[nod.Number]; ok {
				flag--
			}
		}
		if flag == 0 {
			log.Printf("Rule #%d has been prooved", r.Number)
		} else {
			continue
		}

		ds.closedRules[r.Number] = r
		ds.closedRulesOrder = append(ds.closedRulesOrder, r.Number)
		ds.closedNodes[r.Result.Number] = r.Result
		ds.closedNodesOrder = append(ds.closedNodesOrder, r.Result.Number)
		ds.stack.Push(r)

		switch i {
		case len(ds.openRules) - 1:
			ds.openRules = ds.openRules[:i]
		case 0:
			ds.openRules = ds.openRules[1:]
		default:
			ds.openRules = append(ds.openRules[:i], ds.openRules[i+1:]...)
		}

		tar := r.Result
		number := 0
		fl := false
		for i, nod := range ds.openNodes {
			if nod == tar {
				fl = true
				number = i
			}
		}

		if fl {
			switch number {
			case len(ds.openNodes) - 1:
				ds.openNodes = ds.openNodes[:i]
			case 0:
				ds.openNodes = ds.openNodes[1:]
			default:
				ds.openNodes = append(ds.openNodes[:i], ds.openNodes[i+1:]...)
			}
		}

		if r.Result == ds.target {
			return true, true
		} else {
			return false, true
		}
	}

	// old_n := ds.stack.Len()
	// for i, e := range ds.graph.Edges {
	// 	if e.Start == source && e.Label != enums.Closed {
	// 		// если ребро ведёт в запрещённую вершину, помечаем его как закрытое
	// 		if _, ok := ds.forbiddenMap[e.End.Number]; ok {
	// 			ds.graph.Edges[i].Label = enums.Closed
	// 			continue
	// 		}
	// 		// наилучшая ветвь
	// 		ds.stack.Push(e.End)
	// 		if e.End == ds.target {
	// 			return true
	// 		} else {
	// 			return false
	// 		}
	// 	}
	// }
	// // если ни одной вершины не удалось добавить в стек, значит,
	// // у этой вершины нет больше потомков
	// if old_n == ds.stack.Len() {
	// 	v, _ := ds.stack.Pop()
	// 	ds.forbiddenMap[v.Number] = v
	// }
	return false, false
}

// func (ds DeepSearch) unpackStack() []*node.Node {
// 	n := ds.stack.Len()
// 	path := make([]*node.Node, n)
// 	for i := 0; i < n; i++ {
// 		v, _ := ds.stack.Pop()
// 		path[n-i-1] = v
// 	}
// 	return path
// }
