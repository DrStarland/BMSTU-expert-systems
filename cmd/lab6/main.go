package main

import (
	"expert_systems/pkg/algorithms/BFS_logic"
	"expert_systems/pkg/models/logic"
	"log"
)

func main() {
	a := BFS_logic.NewLogicSearch(
		logic.Rule{
			Id: 0,
			Inputs: []logic.Predicate{
				{
					Name: "не умеет летать", Args: []logic.Variable{{Name: "y"}},
				},
				{
					Name: "ловит рыбу под водой", Args: []logic.Variable{{Name: "y"}},
				},
			},
			Result: logic.Predicate{
				Name: "умеет плавать", Args: []logic.Variable{{Name: "y"}},
			},
		},
		logic.Rule{
			Id: 1,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином", Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "ловит рыбу под водой", Args: []logic.Variable{{Name: "x"}},
			},
		},
		logic.Rule{
			Id: 2,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином", Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "является птицей", Args: []logic.Variable{{Name: "x"}},
			},
		},
		logic.Rule{
			Id: 3,
			Inputs: []logic.Predicate{
				{
					Name: "является пингвином", Args: []logic.Variable{{Name: "x"}},
				},
			},
			Result: logic.Predicate{
				Name: "не умеет летать", Args: []logic.Variable{{Name: "x"}},
			},
		},
	)

	log.Printf("Доказываем правило %d: %s", 0, a.Rules[0].String())
	res := a.Prove(
		[]logic.Predicate{
			{
				Name:   "является пингвином",
				Args:   []logic.Variable{{Name: "ПЕН-ПЕН", Const: true}},
				Proved: true,
			},
		},
		0, // индекс целевого правила в базе правил
	)

	log.Printf("Результат: %v\n", res)
}
