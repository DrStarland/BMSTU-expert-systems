package DFS_logic

import (
	formalaparsing "expert_systems/pkg/algorithms/formala-parsing"
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
	"log"
)

func newLogicSituation(facts []types.Runestring, rules []types.Runestring, target types.Runestring) ([]logic.Predicate_r, []logic.Rule_r, logic.Predicate_r) {
	log.Println(facts, len(facts), rules, len(rules), target)

	facts_logic := make([]logic.Predicate_r, 0)
	for _, f := range facts {
		res := formalaparsing.Full_parse_formula(f)
		facts_logic = append(facts_logic, logic.NewPredicate_r(res.(logic.Predicate), 0, nil))
	}

	rules_logic := make([]logic.Rule_r, 0)
	for rule_idx, r := range rules {
		idx := r.IndexOf('→')

		target := r[idx+1:]
		inputs := r[:idx].Split('&')
		target_r := formalaparsing.Full_parse_formula(target)
		inputs_r := make([]any, 0)
		for _, i := range inputs {
			inputs_r = append(inputs_r, formalaparsing.Full_parse_formula(i))
		}
		vars_dict := make(map[string]logic.Var_r)

		_inputs := make([]logic.Predicate_r, 0, len(inputs_r))
		for _, i := range inputs_r {
			_inputs = append(_inputs, logic.NewPredicate_r(i.(logic.Predicate), rule_idx, vars_dict))
		}

		rules_logic = append(rules_logic,
			logic.Rule_r{
				Id:     rule_idx,
				Inputs: _inputs,
				Result: logic.NewPredicate_r(target_r.(logic.Predicate), rule_idx, vars_dict),
			},
		)
	}

	target_logic := logic.NewPredicate_r(formalaparsing.Full_parse_formula(target).(logic.Predicate), len(rules)+1, nil)

	log.Println("Факты:")
	for _, fa := range facts_logic {
		log.Println("\t", fa)
	}
	log.Println("Правила:")
	for _, ruru := range rules_logic {
		log.Println("\t", ruru)
	}
	log.Print("Цель: ", target_logic, "\n\n\n")

	return facts_logic, rules_logic, target_logic
}

func getKeysAsAny[K comparable, V any](x map[K]V) []any {
	keys_x := make([]any, 0, len(x))
	for k := range x {
		keys_x = append(keys_x, k)
	}
	return keys_x
}
