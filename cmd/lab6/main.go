package main

import (
	"expert_systems/pkg/algorithms/BFS_logic"
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
	"log"
)

var Runes = types.Runestring("ᚠ ᚢ ᚦ ᚫ ᚱ ᚲ ᚷ ᚹ ᚺ ᚾ ᛁ ᛃ ᛇ ᛈ ᛉ ᛋ ᛏ ᛒ ᛖ ᛗ ᛚ ᛝ ᛟ ᛞ ᚸ")

func main() {
	a := BFS_logic.NewLogicSearch(
		logic.Rule{
			Id: 0,
			Inputs: []logic.Predicate{
				{
					Name: "не умеет летать",
					Args: []logic.Variable{{Name: "x"}},
				},
				{
					Name: "ловит рыбу под водой",
					Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "умеет плавать",
				Args: []logic.Variable{{Name: "x"}},
			},
		},
		logic.Rule{
			Id: 1,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином",
					Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "ловит рыбу под водой",
				Args: []logic.Variable{{Name: "x"}},
			},
		},
		logic.Rule{
			Id: 2,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином",
					Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "является птицей",
				Args: []logic.Variable{{Name: "x"}},
			},
		},
		logic.Rule{
			Id: 3,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином",
					Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "не умеет летать",
				Args: []logic.Variable{{Name: "x"}},
			},
		},
	)

	res := a.Prove(
		[]logic.Predicate{ // facts =
			{
				Name:   "является пингвином",
				Args:   []logic.Variable{{Name: "x_GERTRUDA", Status: enums.VAL}},
				Proved: true,
			},
			//            logic.Predicate{
			//                name = "не умеет летать",
			//                variables = []logic.Variable{{"X_GERTRUDA", Variable.Status.HAS_VALUE}}
			//            }
		},
		0, // targetRule
	)

	// log.Println("Доказываем правило %d: %s", a.rules)
	log.Println(res)
}
