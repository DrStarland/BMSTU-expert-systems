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
	//// Постоянная память
	// база знаний -- дерево и-или
	knowledgebase and_or_tree.Tree

	//// "рабочая память"
	// целевая вершина
	target *node.Node
	// списки открытых вершин и правил
	openNodes []*node.Node
	openRules []rule.Rule
	// карта закрытых вершин
	closedNodes      map[int]*node.Node
	closedNodesOrder []int // список номеров для хранения порядка
	// карта закрытых правил
	closedRules      map[int]rule.Rule
	closedRulesOrder []int // список номеров для хранения порядка добавления
	// запрещённые карты правил и вершин
	forbiddenRules map[int]rule.Rule
	forbiddenNodes map[int]*node.Node
	// вспомогательный стек для работы алгоритма
	stack StackInterface
}

// "конструктор" алгоритма поиска
func NewSearch(tr and_or_tree.Tree) DeepSearch {
	stck := stack.NewStack[rule.Rule]()
	return DeepSearch{
		knowledgebase: tr,
		target:        nil,
		stack:         stck,
	}
}

/*
Инициализация рабочей памяти алгоритма перед выполнением задачи поиска
*/
func (ds *DeepSearch) init(initial_nodes []*node.Node, target *node.Node) {
	ds.openNodes = []*node.Node{}
	ds.openRules = []rule.Rule{}
	ds.closedNodes = make(map[int]*node.Node, len(initial_nodes))
	ds.closedRules = make(map[int]rule.Rule)
	ds.forbiddenRules = make(map[int]rule.Rule)
	ds.forbiddenNodes = make(map[int]*node.Node)

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
	// очищаем стек (на случай, если это не первый запуск алгоритма)
	for ds.stack.Len() > 0 {
		ds.stack.Pop()
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

	// основной цикл алгоритма
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
	log.Printf("ds.forbiddenNodes: %v\n", ds.forbiddenNodes)
	log.Printf("ds.forbiddenRules: %v\n", ds.forbiddenRules)
	return nil, nil
}

// Проверяет, хватает ли имеющихся узлов (фактов), чтобы доказать правило
func (ds *DeepSearch) checkRuleProvability(r rule.Rule) bool {
	flag := len(r.Inputs)
	// просто перебираем входные узлы правила, проверяя,
	// все ли они есть в списке закрытых вершин
	for _, nod := range r.Inputs {
		if _, ok := ds.closedNodes[nod.Number]; ok {
			flag--
		}
	}
	return flag == 0
}

/*
добавляет переданное правило и его выходную вершину
в карты закрытых правил и вершин соответственно
*/
func (ds *DeepSearch) addRuleAndResultToClosedMaps(r rule.Rule) {
	ds.closedRules[r.Number] = r
	ds.closedRulesOrder = append(ds.closedRulesOrder, r.Number)
	ds.closedNodes[r.Result.Number] = r.Result
	ds.closedNodesOrder = append(ds.closedNodesOrder, r.Result.Number)
}

/*
Ищет позицию переданного правила в списке открытых правил и удаляет оттуда
*/
func (ds *DeepSearch) deleteRuleFromOpenRules(r rule.Rule) {
	for i, v := range ds.openRules {
		if v.Number == r.Number {
			deleteFromArray(&ds.openRules, i)
			break
		}
	}
}

/*
Агрегирующая функция, выполняющая работу по удалению правила и его выходной вершины
из списка открытых правил и вершин
*/
func (ds *DeepSearch) deleteRuleAndResultFromOpenLists(r rule.Rule) {
	ds.deleteResultNodeFromOpenNodes(r.Result)
	ds.deleteRuleFromOpenRules(r)
}

/*
Поиск первого правила, которое можно доказать, и соответствующие этому операции
над рабочей памятью в случае обнаружения.
Синтаксис языка позволяет указать в сигнатуре функции имена переменных,
значения которых нужно вернуть, что позволяет не указывать их явно в return.
По умолчанию значения этих переменных равны false.
*/
func (ds *DeepSearch) findRules() (decisionFlag bool, decisionCanBeFound bool) {
	/* На первой итерации стек пуст. Также он может опустеть в ходе работы алгоритма,
	если рассматриваемая ветвь не может быть доказана.
	Ищем первое правило, ведущее к целевое вершине */
	if ds.stack.Len() == 0 {
		ruleFinded := ds.addToStackRuleResultingToTarget()
		// если стек опустел, а новое правило не было найдено --
		// все возможности решения исчерпаны, решение не может быть найдено
		if !ruleFinded {
			decisionCanBeFound = false
			return
		}
		// иначе работа функции продолжается
	}
	log.Println("Стек не пуст")

	// вытаскиваем из стека правил правило, чтобы исследовать его.
	// в начале итераций это будет первое правило, ведущее к целевой вершине
	r, _ := ds.stack.Peek()
	log.Println(r)
	// проверяем, доказуемо ли текущее рассматриваемое правило
	if ds.checkRuleProvability(r) {
		// если да, то сохраняем его и его выходную вершину
		// в карты закрытых правил и вершин
		ds.addRuleAndResultToClosedMaps(r)
		// удаляем их из списков открытых правил и вершин
		ds.deleteRuleAndResultFromOpenLists(r)
		// если уже получилось доказать целевую вершину -- выходим
		if _, ok := ds.closedNodes[ds.target.Number]; ok {
			decisionFlag = true
			return
		}
		/* даже если в этой ветке не получится найти решение, можно сразу указать,
		   что это возможно, выставив флаг decisionCanBeFound в истинное значение */
		decisionCanBeFound = true
		/* проверяем, можно ли найти решение в таких условиях, если до этого не могли
		   поиск в ширину может найти решение, даже если какие-то правила ещё не были
		   рассмотрены во время обхода в глубину */
		decisionFlag = ds.startBFS()
		/* извлекаем правило из стека, поскольку оно уже в любом случае записано в
		   карту закрытых правил и считается доказанным */
		ds.stack.Pop()
		if decisionFlag {
			// если получилось, опустошаем стек, записывая его содержимое в карту закр. правил
			for ds.stack.Len() > 0 {
				rul, _ := ds.stack.Pop()
				ds.addRuleAndResultToClosedMaps(rul)
			}
			return
		}
	} else {
		// текущее правило не может быть доказано. нужно проверить,
		// возможно ли это в принципе
		log.Println("Правило сходу невозможно доказать")
		// смотрим, можно ли спуститься по дереву ещё ниже. иными словами,
		// есть ли у этого правила входные вершины, которые ещё не попали в список
		// запрещённых, и есть ли правила, которые позволят доказать эти входные вершины
		ancestorFound := ds.findAncestor(r)
		// если ничего не получилось найти, заносим это правило в список запрещённых
		// правил и удаляем из списка открытых
		if !ancestorFound {
			log.Println("Нет потомков")
			// также убираем правило из стека
			ds.stack.Pop()
			ds.forbiddenRules[r.Number] = r
			ds.deleteRuleFromOpenRules(r)
		}
	}
	// если выполнение дошло до этой точки, то ещё можно утверждать, что
	// решение ещё возможно найти на последующих итерациях
	decisionFlag, decisionCanBeFound = false, true
	return
}

/*
используем разработанный ранее обход в ширину, чтобы проверить
*/
func (ds *DeepSearch) startBFS() bool {
	// создаём "объект" алгоритма
	alg := bfs_tree.NewSearch(ds.knowledgebase)
	// в качестве входных вершин задаём закрытые вершины поиска в глубину
	facts := make([]*node.Node, 0, len(ds.closedNodes))
	// переводим карту в массив
	for _, nod := range ds.closedNodes {
		facts = append(facts, nod)
	}

	decisionFlag, _ := alg.FindTarget(ds.target,
		facts...,
	)
	return decisionFlag
}

/*
Поиск предшествующих переданному правилу descendant правил, которые позволяют получить
входные вершины для descendant. Функция возвращает флаг, показывающий, удалось ли
найти что-то и добавить в стек.
*/
func (ds *DeepSearch) findAncestor(descendant rule.Rule) (ancestorFound bool) {
	// просматриваем входные переданного вершины правила
	for _, nod := range descendant.Inputs {
		// проверяем, входит ли данная вершина в закрытые или запрещённые списки
		// ни та, ни другая нас не будет интересовать
		_, nodeIsInClosedNodes := ds.closedNodes[nod.Number]
		_, nodeIsInForbiddenNodes := ds.forbiddenNodes[nod.Number]
		if !nodeIsInForbiddenNodes && !nodeIsInClosedNodes {
			// ищем, выходной вершиной какого правила является эта вершина
			// (ищем первое встречное)
			for _, ruru := range ds.openRules {
				// если это правило не в разделе запрещённых, добавляем его в стек
				_, ok := ds.forbiddenRules[ruru.Number]
				if ruru.Result == nod && !ok {
					ds.stack.Push(ruru)
					// поскольку поиск идёт в глубину, сразу выходим
					return true
				}
			}
			/* если ни одно правило не ведёт к этой вершине, делаем вывод
			что она лежит в самом низу дерева либо ведущие к ней правила являются
			запрещёнными. Следовательно, вершина недосягаема и должна быть запрещена */
			ds.forbiddenNodes[nod.Number] = nod
			ds.deleteResultNodeFromOpenNodes(nod)
			return false
		}
	}

	/* если выполнение функции дошло до этой точки, то в стек не получилось ничего
	добавить. Это правило можно добавить в карту запрещённых, как и его входные
	вершины, не находящиеся в карте закрытых вершин */
	ds.forbiddenRules[descendant.Number] = descendant
	for _, nod := range descendant.Inputs {
		if _, ok := ds.closedNodes[nod.Number]; !ok {
			ds.forbiddenNodes[nod.Number] = nod
			ds.deleteResultNodeFromOpenNodes(nod)
		}
	}

	return false
}

/*
Функция осуществляет поиск первого правила, содержащегося в списке открытых правил,
выходной вершиной которого является целевая вершина.

Возвращаемое значение показывает, удалось ли найти такое правило.
*/
func (ds DeepSearch) addToStackRuleResultingToTarget() (ruleFinded bool) {
	for _, rul := range ds.openRules {
		// проверяем, не содержится ли это правило в списке запрещённых
		_, ok := ds.forbiddenRules[rul.Number]
		if rul.Result == ds.target && !ok {
			ds.stack.Push(rul)
			return true
		}
	}
	return false
}

/*
Функция выполняет поиск позиции целевой вершины в списке (массиве)
открытых вершин и удаляет эту вершину из списка
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
