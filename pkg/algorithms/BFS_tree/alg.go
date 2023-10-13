package bfs_tree

import (
	"expert_systems/pkg/models/and_or_tree"
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/rule"
	"log"
)

type StackInterface interface {
	Len() int
	Push(rule.Rule)
	Pop() (rule.Rule, error)
	Peek() (rule.Rule, error)
}

type BoldSearch struct {
	// Постоянная память
	//// база знаний -- дерево и-или
	knowledgebase and_or_tree.Tree

	// "рабочая память"
	//// целевая вершина
	target           *node.Node
	openNodes        []*node.Node
	openRules        []rule.Rule
	closedNodes      map[int]*node.Node
	closedNodesOrder []int // для хранения порядка добавления
	closedRules      map[int]rule.Rule
	closedRulesOrder []int // для хранения порядка добавления
}

func NewSearch(tr and_or_tree.Tree) BoldSearch {
	// stck := stack.NewStack[rule.Rule]()
	return BoldSearch{
		knowledgebase: tr,
		target:        nil,
		// stack:         stck,
		openNodes:   []*node.Node{},
		openRules:   []rule.Rule{},
		closedNodes: map[int]*node.Node{},
		closedRules: map[int]rule.Rule{},
		// path:          nil,
	}
}

/*
Инициализация рабочей памяти алгоритма перед выполнением задачи поиска
*/
func (ds *BoldSearch) init(initial_nodes []*node.Node, target *node.Node) {
	// назначаем целевую вершину
	ds.target = target
	// добавляем входные вершины в список закрытых
	for _, nod := range initial_nodes {
		ds.closedNodes[nod.Number] = nod
		ds.closedNodesOrder = append(ds.closedNodesOrder, nod.Number)
	}
	// формируем список открытых вершин из базы знаний
	for k, nod := range ds.knowledgebase.Nodes {
		if _, ok := ds.closedNodes[k]; !ok {
			ds.openNodes = append(ds.openNodes, nod)
		}
	}
	// копируем все правила из базы знаний в список открытых правил
	ds.openRules = append(ds.openRules, ds.knowledgebase.Rules...)
	return
}

/*
Основная функция -- поиск целевой вершины
*/
func (ds *BoldSearch) FindTarget(target *node.Node, inputs ...*node.Node) ([]*node.Node, error) {
	// инициализируем "рабочую память"
	ds.init(inputs, target)
	log.Printf("Целевая вершина: №%d", target.Number)
	// флаг, что решение найдено
	decisionFlag := false
	// флаг самой возможности поиска решения. выставляется в ложь, когда
	// по результатам просмотра всей базы открытых правил не удаётся найти
	// хотя бы одно подходящее под текущие данные правило
	decisionCanBeFound := true

	for !decisionFlag && decisionCanBeFound {
		decisionFlag, decisionCanBeFound = ds.findRules()
	}

	if !decisionFlag {
		log.Println("Нет решения")
	} else {
		log.Println("Решение найдено")
	}

	log.Printf(`
	Порядок добавления найденных вершин в процессе решения: %v,
	порядок добавления доказанных правил в процессе решения: %v.`,
		ds.closedNodesOrder, ds.closedRulesOrder,
	)
	return nil, nil
}

// Проверяет, хватает ли имеющихся узлов (фактов), чтобы доказать правило
func (ds *BoldSearch) checkRuleProvability(r rule.Rule) bool {
	flag := len(r.Inputs)
	// просто перебираем входные узлы правила, проверяя, все ли они есть
	// в списке закрытых вершин
	for _, nod := range r.Inputs {
		if _, ok := ds.closedNodes[nod.Number]; ok {
			flag--
		}
	}
	return flag == 0
}

/*
Поиск первого правила, которое можно доказать, и соответствующие этому операции
над рабочей памятью в случае обнаружения
*/
func (ds *BoldSearch) findRules() (bool, bool) {
	// просматриваем базу открытых правил
	for i, r := range ds.openRules {
		flag := ds.checkRuleProvability(r)
		if flag {
			log.Printf("Правило №%d было доказано", r.Number)
			// добавляем это правило и его выходный узел в списки закрытых
			ds.closedRules[r.Number] = r
			ds.closedRulesOrder = append(ds.closedRulesOrder, r.Number)
			ds.closedNodes[r.Result.Number] = r.Result
			ds.closedNodesOrder = append(ds.closedNodesOrder, r.Result.Number)
			// удаляем их из списков открытых правил и вершин
			deleteFromArray(&ds.openRules, i)
			ds.deleteResultNodeFromOpenNodes(r.Result)
			// если была найдена целевая вершина -- сообщаем о нахождении решения,
			// иначе выходим из функции и продолжаем поиск
			if r.Result == ds.target {
				return true, true
			}
			return false, true
		}
	}
	return false, false
}

/*
Функция выполняет поиск позиции целевой вершины в списке (массиве)
открытых вершин и удаляет её
*/
func (ds *BoldSearch) deleteResultNodeFromOpenNodes(target *node.Node) {
	for j, nod := range ds.openNodes {
		// ищем позицию выходной вершины в списке открытых вершин
		if nod == target {
			deleteFromArray(&ds.openNodes, j)
			break
		}
	}
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
