package DFS_logic

import (
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
	"expert_systems/pkg/utils"
	"log"
)

type DeepSearch struct {
	//// Постоянная память
	// база знаний -- дерево и-или
	facts  []logic.Predicate_r
	rules  []logic.Rule_r
	target logic.Predicate_r

	unifier Unifier
}

// "конструктор" алгоритма поиска
func NewSearch(facts []types.Runestring, rules []types.Runestring, target types.Runestring) DeepSearch {
	facts_, rules_, target_ := newLogicSituation(facts, rules, target)
	ds := DeepSearch{
		facts:  facts_,
		rules:  rules_,
		target: target_,
	}
	return ds
}

/*
Инициализация рабочей памяти алгоритма перед выполнением задачи поиска
*/
func (ds *DeepSearch) init() {
	ds.unifier = NewUnifier(ds.facts, ds.rules, ds.target)
}

func (ds *DeepSearch) ProveTarget() bool {
	ds.init()
	proved, _, hm := ds.proveRecursive(ds.target, nil, 0)
	log.Println("Main HM", hm)
	return proved
}

func (ds *DeepSearch) proveRecursive(target logic.Predicate_r, parent any, start_search_idx int) (bool, int, []any) {
	decisionCanBeFound := true

	log.Println(ds.unifier.variables)

	for decisionCanBeFound {
		matched, match_idx := ds.findAtomForUnification(target, start_search_idx)
		if matched == nil {
			log.Println("FORBIDDEN: target ", target, "can not be reached")
			// решение больше не может быть найдено
			return false, match_idx, nil
		}

		switch v := matched.(type) {
		case logic.Predicate_r:
			_, transaction := ds.unifier.unifyPredicate(target, v, false)
			log.Println("PROVED: ", v, v.Vars)
			return true, match_idx, []any{transaction}
		case logic.Rule_r:
			log.Println("PROVING: ", v)
			n := len(v.Inputs)
			transactions := make([][]any, n)
			start_idx := make([]int, n)

			_, base_transaction := ds.unifier.unifyPredicate(target, v.Result, false)

			sub_idx := 0

			for sub_idx < n {
				ok, start_idx_new, trans := ds.proveRecursive(v.Inputs[sub_idx], nil, start_idx[sub_idx])
				if ok {
					start_idx[sub_idx] = start_idx_new + 1
					transactions[sub_idx] = trans
					sub_idx++
				}
				// else {
				// 	transactions[sub_idx] = make([]any, 0)
				// 	start_idx[sub_idx] = 0
				// 	sub_idx--
				// 	if sub_idx < 0 {
				// 		break
				// 	}
				// 	// эта ветка не реализована
				// 	ds.unifier.cancelChanges(transactions[sub_idx])
				// }
			}
			if sub_idx == n {
				_res := []any{base_transaction}

				_res = append(_res, utils.MySum(transactions)...)
				log.Println("TRANSACTIONS", _res)
				return true, match_idx, _res
			} else {
				return false, -1, nil
			}
		}
	}

	return false, -1, nil
}

func (ds *DeepSearch) findAtomForUnification(target logic.Predicate_r, start_search_idx int) (logic.Term, int) {
	for i := start_search_idx; i < len(ds.facts); i++ {
		if ok, _ := ds.unifier.unifyPredicate(target, ds.facts[i], true); ok { //  check_only=true
			return ds.facts[i], i
		}
	}

	for i := start_search_idx; i < len(ds.rules); i++ {
		if ok, _ := ds.unifier.unifyPredicate(target, ds.rules[i].Result, true); ok {
			return ds.rules[i], i
		}
	}
	return nil, -1
}
