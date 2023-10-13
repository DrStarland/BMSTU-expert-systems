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

	log.Println("Правила в базе: ", ds.knowledgebase.Rules)
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
	// флаг самой возможности поиска решения. выставляется в ложь, когда по результатам просмотра всей
	// базы открытых правил не удаётся н
	decisionCanBeFound := true

	for !decisionFlag && decisionCanBeFound {
		decisionFlag, decisionCanBeFound = ds.findRules()
	}

	if !decisionFlag {
		log.Println("Нет решения")
	} else {
		log.Println("Решение найдено")
	}

	log.Printf(
		`
	Порядок добавления найденных вершин в процессе решения: %v,
	порядок добавления доказанных правил в процессе решения: %v.`,
		ds.closedNodesOrder, ds.closedRulesOrder,
	)
	return nil, nil
}

// Проверяет, хватает ли имеющихся узлов (фактов), чтобы доказать правило
func (ds *DeepSearch) checkRuleProvability(r rule.Rule) bool {
	flag := len(r.Inputs)
	for _, nod := range r.Inputs {
		if _, ok := ds.closedNodes[nod.Number]; ok {
			flag--
		}
	}
	return flag == 0
}

func (ds *DeepSearch) findRules() (bool, bool) {
	// просматриваем базу открытых правил
	for i, r := range ds.openRules {
		flag := ds.checkRuleProvability(r)
		if flag {
			log.Printf("Правило №%d было доказано", r.Number)

			ds.closedRules[r.Number] = r
			ds.closedRulesOrder = append(ds.closedRulesOrder, r.Number)
			ds.closedNodes[r.Result.Number] = r.Result
			ds.closedNodesOrder = append(ds.closedNodesOrder, r.Result.Number)
			ds.stack.Push(r)

			deleteFromArray(&ds.openRules, i)

			for j, nod := range ds.openNodes {
				if nod == r.Result {
					deleteFromArray(&ds.openNodes, j)
					break
				}
			}

			if r.Result == ds.target {
				return true, true
			} else {
				return false, true
			}
		}
	}
	return false, false
}

func deleteFromArray(arr interface{}, i int) {
	switch v := arr.(type) {
	case *[]rule.Rule:
		switch i {
		case len(*v) - 1:
			*v = (*v)[:i]
		case 0:
			(*v) = (*v)[1:]
		default:
			(*v) = append((*v)[:i], (*v)[i+1:]...)
		}
	case *[]node.Node:
		switch i {
		case len(*v) - 1:
			*v = (*v)[:i]
		case 0:
			(*v) = (*v)[1:]
		default:
			(*v) = append((*v)[:i], (*v)[i+1:]...)
		}
	}
}
