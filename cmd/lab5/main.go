package main

import (
	resolution "expert_systems/pkg/algorithms/Resolution"
	"expert_systems/pkg/models/logic"
)

const (
	NEG      = true
	POSITIVE = false
)

func main() {
	// P2(x1, y1) | P5(w1) | !P6(z1)
	// P3(C) | !P4(z1) | P1(x1, y1, z1)
	// Доказать !P1(A, B, C)
	f2_1 := logic.Formula{
		Items: []*logic.Disjunct{
			{
				Predicates: []*logic.Predicate{
					{
						Name: "P2", Args: []logic.Variable{logic.NewVariable("x1"), logic.NewVariable("y1")},
					},
					{
						Name: "P5", Args: []logic.Variable{logic.NewVariable("w1")},
					},
					{
						Name: "P6", Negative: true, Args: []logic.Variable{logic.NewVariable("z1")},
					},
				},
			},

			{
				Predicates: []*logic.Predicate{
					{
						Name: "P3", Args: []logic.Variable{logic.NewConst("C")},
					},
					{
						Name: "P4", Negative: true, Args: []logic.Variable{logic.NewVariable("z1")},
					},
					{
						Name: "P1",
						Args: []logic.Variable{logic.NewVariable("x1"), logic.NewVariable("y1"), logic.NewVariable("z1")},
					},
				},
			},
		},
	}

	f2_2 := logic.Formula{
		Items: []*logic.Disjunct{
			{
				Predicates: []*logic.Predicate{
					{
						Name: "P2", Negative: true, Args: []logic.Variable{logic.NewConst("A"), logic.NewConst("B")},
					},
					{
						Name: "P5", Args: []logic.Variable{logic.NewVariable("w2")},
					},
					{
						Name: "P6", Args: []logic.Variable{logic.NewVariable("z2")},
					},
				},
			},

			{
				Predicates: []*logic.Predicate{
					{
						Name: "P4", Args: []logic.Variable{logic.NewVariable("z2")},
					},
					{
						Name: "P3", Negative: true, Args: []logic.Variable{logic.NewVariable("z2")},
					},
				},
			},
		},
	}

	neg_target := logic.Formula{
		Items: []*logic.Disjunct{
			{
				Predicates: []*logic.Predicate{
					{
						Negative: true, Name: "P1",
						Args: []logic.Variable{logic.NewConst("A"), logic.NewConst("B"), logic.NewConst("C")},
					},
				},
			},
		},
	}

	rs := resolution.NewSearch([]logic.Formula{f2_1, f2_2}, neg_target)
	rs.Prove()
}
