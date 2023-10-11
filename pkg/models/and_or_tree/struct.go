package and_or_tree

import (
	"expert_systems/pkg/models/node"
	"expert_systems/pkg/models/rule"

	"github.com/pkg/errors"
)

type Tree struct {
	// основа дерева -- список правил
	Rules []rule.Rule
	// карта всех вершин по номерам
	Nodes map[int]*node.Node
}

type Path struct {
	// само содержание пути
	Data []interface{}
}

// введите структуру графа как список попарных номеров вершин, являющихся началом и концом ребра
// например, NewTree(1, 2, 2, 3, 3, 1)
func appendAndGet(mv map[int]*node.Node, number int) *node.Node {
	v, ok := mv[number]
	if ok {
		return v
	}
	v = node.NewNode(number)
	mv[number] = v
	return v
}

var (
	NotEnoughNodesInRuleError = errors.New("Не хватает вершин, чтобы построить модуль правила")
	EmptyTreeError            = errors.New("Нет правил, чтобы построить дерево")
)

type RuleFormat struct {
	Number        int
	InputsNumbers []int
	ResultNumber  int
}

func NewTree(rules ...RuleFormat) (Tree, error) {
	n := len(rules)
	if n == 0 {
		return Tree{}, EmptyTreeError
	}

	tr := Tree{
		Rules: make([]rule.Rule, n),
		Nodes: make(map[int]*node.Node, n<<2),
	}

	for i, rul := range rules {
		if len(rul.InputsNumbers) < 1 {
			return Tree{}, errors.Wrapf(NotEnoughNodesInRuleError, "rule #%d", i)
		}

		ruleInputNodes := make([]*node.Node, len(rul.InputsNumbers))
		output := appendAndGet(tr.Nodes, rul.ResultNumber)
		for j, number := range rul.InputsNumbers {
			ruleInputNodes[j] = appendAndGet(tr.Nodes, number)
		}
		tr.Rules[i] = rule.NewRule(rul.Number, output, ruleInputNodes...)
	}

	return tr, nil
}
