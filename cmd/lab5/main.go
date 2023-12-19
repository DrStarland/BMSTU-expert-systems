package main

import (
	resolution "expert_systems/pkg/algorithms/Resolution"
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
)

var Runes = types.Runestring("ᚠ ᚢ ᚦ ᚫ ᚱ ᚲ ᚷ ᚹ ᚺ ᚾ ᛁ ᛃ ᛇ ᛈ ᛉ ᛋ ᛏ ᛒ ᛖ ᛗ ᛚ ᛝ ᛟ ᛞ ᚸ")

const (
	NEG = true
	POS = false
)

func main() {
	// P2(x1, y1) | P5(w1) | !P6(z1)
	// P3(C) | !P4(z1) | P1(x1, y1, z1)
	f2_1 := logic.Formula{
		Items: []*logic.Disjunct{
			{
				Predicates: []*logic.Predicate{
					{
						Name:     "P2",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("x1"),
							logic.NewVariable("y1"),
						},
					},
					{
						Name:     "P5",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("w1"),
						},
					},
					{
						Name:     "P6",
						Negative: NEG,
						Args: []logic.Variable{
							logic.NewVariable("z1"),
						},
					},
				},
			},

			{
				Predicates: []*logic.Predicate{
					{
						Name:     "P3",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewConst("C"),
						},
					},
					{
						Name:     "P4",
						Negative: NEG,
						Args: []logic.Variable{
							logic.NewVariable("z1"),
						},
					},
					{
						Name:     "P1",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("x1"),
							logic.NewVariable("y1"),
							logic.NewVariable("z1"),
						},
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
						Name:     "P2",
						Negative: NEG,
						Args: []logic.Variable{
							logic.NewConst("A"),
							logic.NewConst("B"),
						},
					},
					{
						Name:     "P5",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("w2"),
						},
					},
					{
						Name:     "P6",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("z2"),
						},
					},
				},
			},

			{
				Predicates: []*logic.Predicate{
					{
						Name:     "P4",
						Negative: POS,
						Args: []logic.Variable{
							logic.NewVariable("z2"),
						},
					},
					{
						Name:     "P3",
						Negative: NEG,
						Args: []logic.Variable{
							logic.NewVariable("z2"),
						},
					},
				},
			},
		},
	}

	// // _ !P1(A, B, C)
	neg_target := logic.Formula{
		Items: []*logic.Disjunct{
			{
				Predicates: []*logic.Predicate{
					{
						Name:     "P1",
						Negative: NEG,
						Args: []logic.Variable{
							logic.NewConst("A"),
							logic.NewConst("B"),
							logic.NewConst("C"),
						},
					},
				},
				// logic.Predicate { "P1", NEG, { Term::const_("A"), Term::const_("B"), Term::const_("C") } },
			},
		},
	}

	rs := resolution.NewSearch([]logic.Formula{f2_1, f2_2}, neg_target)
	rs.Solve()

}
