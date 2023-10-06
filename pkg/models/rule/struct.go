package rule

import "expert_systems/pkg/models/node"

type Rule struct {
	Number  int
	Feature int

	Inputs []*node.Node
	Result *node.Node
}

func NewRule(number int, result *node.Node, inputs ...*node.Node) *Rule {
	return &Rule{
		Number:  number,
		Feature: 0,

		Inputs: inputs,
		Result: result,
	}
}
