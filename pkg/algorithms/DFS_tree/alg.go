package DFS_tree

import (
	bfs_tree "expert_systems/pkg/algorithms/BFS_tree"
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
	forbiddenRules   map[int]rule.Rule
	forbiddenNodes   map[int]*node.Node
	stack            StackInterface
}

func NewSearch(tr and_or_tree.Tree) DeepSearch {
	stck := stack.NewStack[rule.Rule]()
	return DeepSearch{
		knowledgebase:  tr,
		target:         nil,
		stack:          stck,
		openNodes:      []*node.Node{},
		openRules:      []rule.Rule{},
		closedNodes:    map[int]*node.Node{},
		closedRules:    map[int]rule.Rule{},
		forbiddenRules: make(map[int]rule.Rule),
		forbiddenNodes: make(map[int]*node.Node),
		// path:          nil,
	}
}

/*
Инициализация рабочей памяти алгоритма перед выполнением задачи поиска
*/
func (ds *DeepSearch) init(initial_nodes []*node.Node, target *node.Node) {
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
	// очищаем стек
	for ds.stack.Len() > 0 {
		ds.stack.Pop()
	}
	// ищем правило, которое порождает целевую вершину, чтобы добавить его в стек первым
	for _, ruru := range ds.openRules {
		if ruru.Result == target {
			ds.stack.Push(ruru)
			// deleteFromArray(&ds.openRules, i)
			break
		}
	}
	return
}

/*
Основная функция -- поиск целевой вершины
*/
func (ds *DeepSearch) FindTarget(target *node.Node, inputs ...*node.Node) ([]*node.Node, error) {
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

	// for i, j := 0, len(ds.closedRulesOrder)-1; i < j; i, j = i+1, j-1 {
	// 	ds.closedRulesOrder[i], ds.closedRulesOrder[j] = ds.closedRulesOrder[j], ds.closedRulesOrder[i]
	// }

	log.Printf(`
	Порядок добавления найденных вершин в процессе решения: %v,
	порядок добавления доказанных правил в процессе решения: %v.`,
		ds.closedNodesOrder, ds.closedRulesOrder,
	)
	log.Printf("ds.forbiddenNodes: %v\n", ds.forbiddenNodes)
	log.Printf("ds.forbiddenRules: %v\n", ds.forbiddenRules)
	return nil, nil
}

// Проверяет, хватает ли имеющихся узлов (фактов), чтобы доказать правило
func (ds *DeepSearch) checkRuleProvability(r rule.Rule) bool {
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
func (ds *DeepSearch) findRules() (decisionFlag bool, decisionCanBeFound bool) {
	// если стек опустел -- все возможности решения исчерпаны,
	// решение не может быть найдено
	if ds.stack.Len() == 0 {
		decisionFlag, decisionCanBeFound = false, false
		return
	}
	log.Println("Стек не пуст")

	// вытаскиваем из стека правил правило, чтобы исследовать его.
	// в начале итераций это будет первое правило, ведущее к целевой вершине
	r, _ := ds.stack.Peek()
	log.Println(r)
	// просматриваем базу открытых правил
	if ds.checkRuleProvability(r) {
		ds.closedRules[r.Number] = r
		ds.closedRulesOrder = append(ds.closedRulesOrder, r.Number)
		ds.closedNodes[r.Result.Number] = r.Result
		ds.closedNodesOrder = append(ds.closedNodesOrder, r.Result.Number)

		// удаляем их из списков открытых правил и вершин
		ds.deleteResultNodeFromOpenNodes(r.Result)
		for i, v := range ds.openRules {
			if v.Number == r.Number {
				deleteFromArray(&ds.openRules, i)
				break
			}
		}
		// если уже получилось доказать целевую вершину -- выходим
		if _, ok := ds.closedNodes[ds.target.Number]; ok {
			decisionFlag, decisionCanBeFound = true, true
			return
		}
		decisionCanBeFound = true
		// проверяем, можно ли найти решение в таких условиях, если до этого не могли
		decisionFlag = ds.startBFS()
		log.Println("После бек-трекинга")
		ds.stack.Pop()
		if decisionFlag {
			// если получилось, опустошаем стек
			for ds.stack.Len() > 0 {
				ruru, _ := ds.stack.Pop()
				ds.closedRules[ruru.Number] = ruru
				ds.closedRulesOrder = append(ds.closedRulesOrder, ruru.Number)
				ds.closedNodes[ruru.Result.Number] = ruru.Result
				ds.closedNodesOrder = append(ds.closedNodesOrder, ruru.Result.Number)
			}
			return
		}
		log.Println("Я здесь 2")
		ancestorFound := ds.findAncestor(r)
		if !ancestorFound {
			ds.stack.Pop()
			ds.forbiddenRules[r.Number] = r

			for i, v := range ds.openRules {
				if v.Number == r.Number {
					deleteFromArray(&ds.openRules, i)
					break
				}
			}
		}
		return
	} else {
		log.Println("Правило сходу невозможно доказать")
		ancestorFound := ds.findAncestor(r)
		if !ancestorFound {
			log.Println("Нет потомков")
			ds.stack.Pop()
			ds.forbiddenRules[r.Number] = r

			for i, v := range ds.openRules {
				if v.Number == r.Number {
					deleteFromArray(&ds.openRules, i)
					break
				}
			}
		}
	}

	decisionFlag, decisionCanBeFound = false, true
	return
}

func (ds *DeepSearch) startBFS() bool {
	alg := bfs_tree.NewSearch(ds.knowledgebase)
	facts := make([]*node.Node, 0, len(ds.closedNodes))
	for _, nod := range ds.closedNodes {
		facts = append(facts, nod)
	}

	log.Println(facts)

	decisionFlag, _ := alg.FindTarget(ds.target,
		facts...,
	)
	return decisionFlag
}

/*
 */
func (ds *DeepSearch) findAncestor(descendant rule.Rule) (ancestorFound bool) {
	// просматриваем входные вершины правила
	for _, nod := range descendant.Inputs {
		// проверяем, входит ли данная вершина в закрытые или запрещённые списки
		// ни та, ни другая нас не будет интересовать
		_, nodeInClosedNodes := ds.closedNodes[nod.Number]
		_, nodeInForbiddenNodes := ds.forbiddenNodes[nod.Number]
		if !nodeInForbiddenNodes && !nodeInClosedNodes {
			// ищем, выходной вершиной какого правила является эта вершина
			// (ищем первое встречное)
			for _, ruru := range ds.openRules {
				// если это правило не в разделе запрещённых, добавляем
				// его в стек
				_, ok := ds.forbiddenRules[ruru.Number]
				if ruru.Result == nod && !ok {
					ds.stack.Push(ruru)
					return true
				}
			}
			// если ни одно правило не ведёт к этой вершине, делаем вывод
			// что она лежит в самом низу дерева и, следовательно, недосягаема,
			// если не была задана изначально
			ds.forbiddenNodes[nod.Number] = nod
			ds.deleteResultNodeFromOpenNodes(nod)
			return false
		}
	}

	ds.forbiddenRules[descendant.Number] = descendant
	for _, nod := range descendant.Inputs {
		ds.forbiddenNodes[nod.Number] = nod
		ds.deleteResultNodeFromOpenNodes(nod)
	}

	return false
}

/*
Функция выполняет поиск позиции целевой вершины в списке (массиве)
открытых вершин и удаляет её
*/
func (ds *DeepSearch) deleteResultNodeFromOpenNodes(target *node.Node) {
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

// class AlgorithmDFS:
//     def __init__(self, graph: models.Graph):
//         self._graph: models.Graph = graph

//     def search(
//         self,
//         in_vertexes: models.Vertex,
//         target: models.Vertex
//     ):
//         for vertex in self._graph.vertexes():
//             if vertex.number in in_vertexes:
//                 vertex.state = models.VertexState.APPROVED

//         stack: Deque = deque()
//         forbidden_rules: List[models.Rule] = []
//         approved_rules: List[models.Rule] = []

//         stack.append(target)
//         while stack:
//             value = stack[-1]
//             # цель вершина
//             if (isinstance(value, models.Vertex)):
//                 target_vertex = typing.cast(models.Vertex, value)

//                 # если уже доказано
//                 if target_vertex.state == models.VertexState.APPROVED:
//                     stack.pop()
//                     continue

//                 # правило которое доказывает вершину
//                 target_rule = self._search_rule(target_vertex)
//                 if not target_rule:
//                     target_vertex.state = models.VertexState.FORBIDDEN
//                     stack.pop()
//                 elif target_rule.state == models.RuleState.APPROVED:
//                     target_vertex.state = models.VertexState.APPROVED
//                 else:
//                     stack.append(target_rule)
//             # цель правило
//             else:
//                 target_rule = typing.cast(models.Rule, value)
//                 # вершина которая не доказана
//                 target_vertex = self._search_vertex(target_rule)

//                 # нет недоказанных вершин
//                 if not target_vertex:
//                     target_rule.state = models.RuleState.APPROVED
//                     approved_rules.append(target_rule)
//                     stack.pop()

//                 # входная вершина запрещена
//                 elif target_vertex.state == models.VertexState.FORBIDDEN:
//                     target_rule.state = models.RuleState.FORBIDDEN
//                     forbidden_rules.append(target_rule)
//                     stack.pop()

//                 # иначе вершина неизвестна
//                 else:
//                     stack.append(target_vertex)

//         return approved_rules, forbidden_rules

//     def _search_rule(
//         self,
//         target: models.Vertex,
//     ) -> Optional[models.Rule]:
//         for rule in self._graph.rules():
//             if rule.out_vertex == target \
//                 and rule.state != models.RuleState.FORBIDDEN:
//                 return rule
//         return None

//     def _search_vertex(
//         self,
//         rule: models.Rule,
//     ) -> Optional[models.Vertex]:
//         for vertex in rule.in_vertexes:
//             if vertex.state != models.VertexState.APPROVED:
//                 return vertex
//         return None
