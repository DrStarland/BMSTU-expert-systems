package BFS_logic

import (
	"expert_systems/pkg/models/logic"
	"log"
)

type LogicSearch struct {
	// база знаний
	rules []logic.Rule
	// рабочая память
	closedRules set[int]
	closedFacts []logic.Predicate
}

type set[T comparable] struct {
	list []T
	map_ map[T]bool
}

func NewSet[T comparable]() set[T] {
	return set[T]{
		list: make([]T, 0),
		map_: map[T]bool{},
	}
}

func (s *set[T]) Add(elems ...T) {
	for _, x := range elems {
		// если такого значения нет в карте, карта в го возвращает стандартное значение типа:
		// в данном случае это будет false
		if ok := s.map_[x]; !ok {
			s.list = append(s.list, x)
			s.map_[x] = true
		}
	}
}

func NewLogicSearch(rules ...logic.Rule) LogicSearch {
	return LogicSearch{
		rules:       rules,
		closedRules: NewSet[int](),
		closedFacts: []logic.Predicate{},
	}
}

func (s set[T]) Size() int {
	return len(s.list)
}

func (s set[T]) Contains(target T) bool {
	_, ok := s.map_[target]
	return ok
}

func (ls *LogicSearch) Prove(facts []logic.Predicate, targetRule int) bool {
	ls.closedFacts = append(ls.closedFacts, facts...)

	var decisionCanBeFound = true
	for decisionCanBeFound && !ls.closedRules.Contains(targetRule) {
		oldSize := len(ls.closedFacts)
		ls.findRules()
		decisionCanBeFound = len(ls.closedFacts) != oldSize
	}

	log.Println("___________________________________________")
	log.Println(len(ls.closedRules.list), ls.closedRules.list)
	log.Println("Facts", ls.closedFacts)

	return ls.closedRules.Contains(targetRule)
}

func (ls *LogicSearch) findRules() {
	for k, rule := range ls.rules {
		if !rule.Proved {
			log.Printf("*** Рассматриваем правило %v\n", rule)
			for i, pred := range rule.Inputs {
				if !pred.Proved {
					for _, it := range ls.closedFacts {
						if pred.Name == it.Name && len(pred.Args) == len(it.Args) {
							log.Printf("Рассматриваем атом `%v` правила %d\n", pred, rule.Id)
							ls.rules[k].Inputs[i].Args = it.Args
							ls.rules[k].Inputs[i].Proved = true
							log.Printf("Доказан атом `%v`\n", pred)
							break
						}
					}
				} else {
					log.Printf("Атом `%v` правила %v уже доказан\n", pred, rule.Id)
				}
			}

			if ls.checkRuleProvability(rule) {
				log.Printf("Правило `%d` (%s) доказано!", rule.Id, rule.String())
				for j, resVar := range rule.Result.Args {
					all_vars := make([]logic.Variable, 0)
					for _, em := range rule.Inputs {
						all_vars = append(all_vars, em.Args...)
					}
					new_name := ""
					for _, it := range all_vars {
						if it.Name[0] == resVar.Name[0] {
							new_name = it.Name
							break
						}
					}
					ls.rules[k].Result.Args[j].Name = new_name
				}
				ls.rules[k].Proved = true
				log.Printf("Добавляем атом `%v`", rule.Result)
				ls.closedRules.Add(rule.Id)
				ls.closedFacts = append(ls.closedFacts, rule.Result)
			}
		}
	}
	log.Println()
	log.Println("В результате итерации список закрытых фактов: ", ls.closedFacts)
	log.Println("=====================")
}

// Проверяет, хватает ли имеющихся узлов (фактов), чтобы доказать правило
func (ls *LogicSearch) checkRuleProvability(rule logic.Rule) bool {
	flag := true
	for _, hm := range rule.Inputs {
		if !hm.Proved {
			flag = false
			break
		}
	}
	return flag
}

// aa := logic.Predicate{Name: "Kek", Args: []logic.Variable{{Name: "LOL", Const: false}}}
// bb := logic.Predicate{Name: "Kek", Args: []logic.Variable{{Name: "Cheburek", Const: false}}}
// cc := logic.Predicate{Name: "Kek", Args: []logic.Variable{{Name: "Shrek", Const: true}}}
// kb := []logic.Predicate{aa, bb, cc}
// log.Println(Unify(kb, &aa, &bb), "LB", kb, "AA", aa, "BB", bb)
// log.Println(Unify(kb, &bb, &cc), "LB", kb, "AA", aa, "BB", bb, "CC", cc)

// func Unify(kb []logic.Predicate, a1, a2 *logic.Predicate) bool {
// 	type pair struct {
// 		first  string
// 		second string
// 	}

// 	if a1.Name != a2.Name {
// 		log.Panic(a1.Name, a2.Name)
// 	}

// 	// Предварительные списки замен
// 	constsmap_pings := make([]pair, 0)
// 	linked_vars := make([]pair, 0)

// 	// Сопоставление аргументов предикатов
// 	for i := 0; i < len(a1.Args); i++ {
// 		arg1, arg2 := &a1.Args[i], &a2.Args[i]
// 		switch {
// 		case arg1.Const && arg2.Const: // Обе константы
// 			if arg1.Name != arg2.Name {
// 				return false
// 			}
// 		case !arg1.Const && !arg2.Const: // Обе переменные
// 			if arg1.Name != arg2.Name {
// 				linked_vars = append(linked_vars, pair{arg1.Name, arg2.Name})
// 			}
// 		default: // Константа и переменная
// 			if arg1.Const {
// 				a1.Args[i], a2.Args[i] = a2.Args[i], a1.Args[i]
// 			} // arg1 - var, arg2 - const
// 			constsmap_pings = append(constsmap_pings, pair{arg1.Name, arg2.Name})
// 		}
// 	}

// 	// Объединение связанных переменных
// 	counter := 1

// 	new_vars := make(map[string]int, 0)
// 	for _, tuple := range linked_vars {
// 		var1, var2 := tuple.first, tuple.second
// 		num1, ok := new_vars[var1]
// 		num2, ok2 := new_vars[var2]
// 		if ok && ok2 {
// 			for vari, num := range new_vars {
// 				if num == num2 {
// 					new_vars[vari] = num1
// 				}
// 			}
// 		} else {
// 			if vari, ok := new_vars[var1]; ok {
// 				new_vars[var2] = vari
// 			} else if vari, ok := new_vars[var2]; ok {
// 				new_vars[var1] = vari
// 			} else {
// 				new_vars[var1] = counter
// 				new_vars[var2] = counter
// 				counter++
// 			}
// 		}
// 	}

// 	// Применение связанных переменных к списку замен констант
// 	for vari, num := range new_vars {
// 		for i, para := range constsmap_pings {
// 			old_v, _ := para.first, para.second
// 			if old_v == vari {
// 				constsmap_pings[i].first = fmt.Sprintf("@%d", num)
// 			}
// 		}
// 	}

// 	// Проверка замен констант на возможность унификации (нет двух разных замен одной переменной)
// 	vars_vals := map[string]string{}
// 	for _, para := range constsmap_pings {
// 		old_v, new_v := para.first, para.second
// 		if v, ok := vars_vals[old_v]; ok {
// 			if v != new_v {
// 				return false
// 			}
// 		} else {
// 			vars_vals[old_v] = new_v
// 		}
// 	}

// 	// Замена связанных переменных
// 	for vari, num := range new_vars {
// 		new_name := fmt.Sprintf("@%d", num)
// 		substitute(kb, vari, new_name, false)
// 	}
// 	// Замена констант
// 	for old_v, new_v := range vars_vals {
// 		substitute(kb, old_v, new_v, true)
// 	}

// 	return true
// }

// func substitute(kb []logic.Predicate, old_v, new_v string, make_const bool) {
// 	_const := ""
// 	if make_const {
// 		_const = " const"
// 	}
// 	log.Println("  Замена: " + old_v + " -> " + new_v + _const)

// 	for j := range kb {
// 		for k, vari := range kb[j].Args {
// 			if vari.Name == old_v {
// 				vari.Name, vari.Const = new_v, make_const
// 				kb[j].Args[k] = vari
// 			}
// 		}
// 	}
// }
