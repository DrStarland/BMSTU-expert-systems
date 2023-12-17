package resolution

import (
	"expert_systems/pkg/models/logic"
	"fmt"
	"log"
)

// Проверка пары дизъюнктов на наличие контрактной пары, и возможная унификация
func (rs *ResolutionSolver) check_disjuncts(d1, d2 *logic.Disjunct) bool {
	for i1 := 0; i1 < len(d1.Predicates); i1++ {
		for i2 := 0; i2 < len(d2.Predicates); i2++ {
			Predicate1 := d1.Predicates[i1]
			Predicate2 := d2.Predicates[i2]

			if !(Predicate1.Name == Predicate2.Name && Predicate1.Negative != Predicate2.Negative) {
				continue
			}
			if rs.new_disjunct_present(d1, i1) && rs.new_disjunct_present(d2, i2) {
				continue
			}

			log.Printf("Унификация: %s и %s", Predicate1.String(), Predicate2.String())
			unified := rs.unify_Predicates(Predicate1, Predicate2)
			if !unified {
				log.Println(" невозможна")
				continue
			}
			log.Println()

			log.Println("Новые:")
			rs.add_new_disjunct(d1, i1)
			rs.add_new_disjunct(d2, i2)
			return true
		}
	}
	return false
}

// Заменить все вхождения переменной old_v на переменную или константу nev_v
func (rs *ResolutionSolver) substitute_Variables(old_v, new_v string, make_const bool) {
	_const := ""
	if make_const {
		_const = " const"
	}
	log.Println("  Замена: " + old_v + " -> " + new_v + _const)
	for i := range rs.disjuncts {
		for j := range rs.disjuncts[i].Predicates {
			for k, vari := range rs.disjuncts[i].Predicates[j].Args {
				if vari.Name == old_v {
					// assert(!Variable.is_const);
					vari.Name = new_v
					if make_const {
						vari.Const = true
					}
					rs.disjuncts[i].Predicates[j].Args[k] = vari
				}
			}
		}
	}
}

// Сформировать дизъюнк из base путём исключения атома с индексом out_idx
func (rs *ResolutionSolver) get_new_disjunct(base *logic.Disjunct, out_idx int) (logic.Disjunct, bool) {
	new_disjunct := logic.NewDisjunct(make([]*logic.Predicate, 0, len(base.Predicates)-1))
	for j := 0; j < len(base.Predicates); j++ {
		if j != out_idx {
			new_disjunct.Predicates = append(new_disjunct.Predicates, base.Predicates[j])
		}
	}

	present := false
	for _, d := range rs.disjuncts {
		if new_disjunct.EqualTo(*d) {
			present = true
			break
		}
	}

	return new_disjunct, present
}

// Есть ли заданный дизъюнкт в списке
func (rs *ResolutionSolver) new_disjunct_present(base *logic.Disjunct, out_idx int) bool {
	_, flag := rs.get_new_disjunct(base, out_idx)
	return flag
}

// Добавить дизъюнкт в список
func (rs *ResolutionSolver) add_new_disjunct(base *logic.Disjunct, out_idx int) bool {
	// log.Println("место добавления: ", rs.disjuncts)
	disj, present := rs.get_new_disjunct(base, out_idx)
	log.Println("страх и ужас", disj)
	if present {
		return false
	}

	log.Println("  " + disj.String())
	if len(disj.Predicates) == 0 {
		rs.final_result = true
	}

	rs.disjuncts = append(rs.disjuncts, &disj)
	// log.Println("место добавления: ", rs.disjuncts)
	// panic("kek")
	return true
}

// Попытаться унифицировать 2 атома, возвращает true, если унифицировано
func (rs *ResolutionSolver) unify_Predicates(a1, a2 *logic.Predicate) bool {
	type pair struct {
		first  string
		second string
	}

	if a1.Name != a2.Name {
		log.Panic(a1.Name, a2.Name)
	}
	if a1.Negative == a2.Negative {
		log.Panic(a1.Negative, a2.Negative)
	}

	// Предварительные списки замен
	consts_mappings := make([]pair, 0)
	linked_vars := make([]pair, 0)

	// Сопоставление аргументов предикатов
	for i := 0; i < len(a1.Args); i++ {
		arg1 := &a1.Args[i]
		arg2 := &a2.Args[i]

		if arg1.Const && arg2.Const { // Обе константы
			if arg1.Name != arg2.Name {
				return false
			}
		} else if !arg1.Const && !arg2.Const { // Обе переменные
			if arg1.Name != arg2.Name {
				linked_vars = append(linked_vars, pair{arg1.Name, arg2.Name})
			}
		} else { // Константа и переменная
			if arg1.Const {
				a1.Args[i], a2.Args[i] = a2.Args[i], a1.Args[i]
				// *arg1, *arg2 = *arg2, *arg1
			} // arg1 - var, arg2 - const
			consts_mappings = append(consts_mappings, pair{arg1.Name, arg2.Name})
		}
	}

	log.Println("LINKED VARS", linked_vars)
	log.Println("СЩТЫЕ ЬФЗЗШТПЫ", consts_mappings)
	// panic("lol")

	// Объединение связанных переменных
	counter := 1

	new_vars := make(map[string]int, 0)
	for _, tuple := range linked_vars {
		var1, var2 := tuple.first, tuple.second
		num1, ok := new_vars[var1]
		num2, ok2 := new_vars[var2]

		if ok && ok2 {
			for vari, num := range new_vars {
				if num == num2 {
					new_vars[vari] = num1
				}
			}
		} else {
			if hm, ok := new_vars[var1]; ok {
				new_vars[var2] = hm
			} else if hm, ok := new_vars[var2]; ok {
				new_vars[var1] = hm
			} else {
				new_num := counter
				counter++
				new_vars[var1] = new_num
				new_vars[var2] = new_num
			}
		}
	}

	log.Println("Counter", counter)
	// if counter > 1 {
	// 	panic("AAA")
	// }
	// Применение связанных переменных к списку замен констант
	for vari, num := range new_vars {
		//  for (auto &[old_v, new_v] : consts_mappings)
		for i, para := range consts_mappings {
			old_v, _ := para.first, para.second
			if old_v == vari {
				consts_mappings[i].first = fmt.Sprintf("@%d", num)
			}
		}
	}

	// // Проверка замен констант на возможность унификации (нет двух разных замен одной переменной)
	// std::map<std::string, std::string> vars_vals;
	// for (auto &[old_v, new_v] : consts_mappings) {
	// 	if (vars_vals.contains(old_v)) {
	// 		if (vars_vals.at(old_v) != new_v)
	// 			return false;
	// 	}
	// 	else {
	// 		vars_vals.emplace(old_v, new_v);
	// 	}
	// }
	// Проверка замен констант на возможность унификации (нет двух разных замен одной переменной)
	vars_vals := map[string]string{}
	for _, para := range consts_mappings {
		old_v, new_v := para.first, para.second
		if v, ok := vars_vals[old_v]; ok {
			if v != new_v {
				return false
			}
		} else {
			vars_vals[old_v] = new_v
		}
	}

	// Замена связанных переменных
	for vari, num := range new_vars {
		new_name := fmt.Sprintf("@%d", num)
		rs.substitute_Variables(vari, new_name, false)
	}
	// Замена констант
	for old_v, new_v := range vars_vals {
		rs.substitute_Variables(old_v, new_v, true)
	}
	log.Println("END: ", new_vars, vars_vals)
	// if death == 3 {
	// 	panic(counter)
	// }
	death++
	return true
}

var death int = 1

func (rs *ResolutionSolver) print_disjuncts() {
	log.Println("Дизъюнкты:\n")
	for _, d := range rs.disjuncts {
		log.Printf("  %s\n", d)
	}

	log.Println("=======================================")
}
