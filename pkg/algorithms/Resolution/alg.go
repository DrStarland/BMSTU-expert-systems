package resolution

import (
	"expert_systems/pkg/models/logic"
	"fmt"
	"log"
)

type ResolutionSearch struct {
	disjuncts []*logic.Disjunct
}

// "конструктор" алгоритма поиска
func NewSearch(formulas []logic.Formula, neg_target logic.Formula) ResolutionSearch {
	rs := ResolutionSearch{
		disjuncts: make([]*logic.Disjunct, 0, len(formulas)+len(neg_target.Items)),
	}

	for _, f := range formulas {
		rs.disjuncts = append(rs.disjuncts, f.Items...)
	}
	rs.disjuncts = append(rs.disjuncts, neg_target.Items...)
	return rs
}

func (rs *ResolutionSearch) Prove() {
	log.Println("Перед началом решения:")
	rs.print_disjuncts()
	decisionCanBeFound := true
	decisionFound := false

	// Основной цикл обхода по базе знаний
	for decisionCanBeFound && !decisionFound {
	iteration:
		for i, d1 := range rs.disjuncts {
			for j, d2 := range rs.disjuncts {
				if d1 != d2 {
					// если дизъюнкты не равны, то смотрить, есть ли в них контрарные пары
					decisionCanBeFound, decisionFound = rs.findOppositePair(rs.disjuncts[i], rs.disjuncts[j])
					if decisionCanBeFound {
						// если пара нашлась, то надо начинать итерации сначала
						// (кроме того, решение уже могло быть найдено)
						break iteration
					}
				}
			}
		}
		rs.print_disjuncts()
	}

	if decisionFound {
		log.Println("Доказано")
	} else {
		log.Println("He доказано")
	}
}

// Проверка пары дизъюнктов на наличие контрарной пары и возможная унификация
func (rs *ResolutionSearch) findOppositePair(d1, d2 *logic.Disjunct) (decisionCanBeFound, decisionFound bool) {
	for i := 0; i < len(d1.Predicates); i++ {
		for j := 0; j < len(d2.Predicates); j++ {
			pred1 := d1.Predicates[i]
			pred2 := d2.Predicates[j]
			if pred1.Name != pred2.Name || pred1.Negative == pred2.Negative {
				continue
			}
			if rs.new_disjunct_present(d1, i) && rs.new_disjunct_present(d2, j) {
				continue
			}

			log.Printf("Унификация: %s и %s", pred1.String(), pred2.String())
			if !rs.Unify(pred1, pred2) {
				log.Println(" не удалась")
				continue
			}

			log.Println("Новые:")
			if rs.add_new_disjunct(d1, i) || rs.add_new_disjunct(d2, j) {
				return true, true
			}
			return true, false
		}
	}
	return false, false
}

// Заменить все вхождения переменной old_vari на переменную или константу nev_v
func (rs *ResolutionSearch) substitute(old_vari, new_vari string, make_const bool) {
	_const := ""
	if make_const {
		_const = " const"
	}
	log.Println("  Замена: " + old_vari + " -> " + new_vari + _const)

	for i := range rs.disjuncts {
		for j := range rs.disjuncts[i].Predicates {
			for k, vari := range rs.disjuncts[i].Predicates[j].Args {
				if vari.Name == old_vari {
					vari.Name, vari.Const = new_vari, make_const
					rs.disjuncts[i].Predicates[j].Args[k] = vari
				}
			}
		}
	}
}

// Сформировать новый дизъюнкт путём исключения предиката по индексу
func (rs *ResolutionSearch) get_new_disjunct(old *logic.Disjunct, out_idx int) (logic.Disjunct, bool) {
	new_disjunct := logic.NewDisjunct(make([]*logic.Predicate, 0, len(old.Predicates)-1))
	for j := 0; j < len(old.Predicates); j++ {
		if j != out_idx {
			new_disjunct.Predicates = append(new_disjunct.Predicates, old.Predicates[j])
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
func (rs *ResolutionSearch) new_disjunct_present(old *logic.Disjunct, out_idx int) bool {
	_, flag := rs.get_new_disjunct(old, out_idx)
	return flag
}

// Добавить дизъюнкт в список
func (rs *ResolutionSearch) add_new_disjunct(old *logic.Disjunct, out_idx int) (decisionFound bool) {
	disj, present := rs.get_new_disjunct(old, out_idx)
	if present {
		return false
	}

	log.Println("  " + disj.String())

	rs.disjuncts = append(rs.disjuncts, &disj)
	// проверяем, появился ли в базе пустой дизъюнкт
	return len(disj.Predicates) == 0
}

// попытка унификации 2 предикатов
func (rs *ResolutionSearch) Unify(a1, a2 *logic.Predicate) bool {
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
		arg1, arg2 := &a1.Args[i], &a2.Args[i]
		switch {
		case arg1.Const && arg2.Const: // Обе константы
			if arg1.Name != arg2.Name {
				return false
			}
		case !arg1.Const && !arg2.Const: // Обе переменные
			if arg1.Name != arg2.Name {
				linked_vars = append(linked_vars, pair{arg1.Name, arg2.Name})
			}
		default: // Константа и переменная
			if arg1.Const {
				a1.Args[i], a2.Args[i] = a2.Args[i], a1.Args[i]
			} // arg1 - var, arg2 - const
			consts_mappings = append(consts_mappings, pair{arg1.Name, arg2.Name})
		}
	}

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
			if vari, ok := new_vars[var1]; ok {
				new_vars[var2] = vari
			} else if vari, ok := new_vars[var2]; ok {
				new_vars[var1] = vari
			} else {
				new_vars[var1] = counter
				new_vars[var2] = counter
				counter++
			}
		}
	}

	// Применение связанных переменных к списку замен констант
	for vari, num := range new_vars {
		for i, para := range consts_mappings {
			old_vari, _ := para.first, para.second
			if old_vari == vari {
				consts_mappings[i].first = fmt.Sprintf("&%d", num)
			}
		}
	}

	// Проверка замен констант на возможность унификации (нет двух разных замен одной переменной)
	vars_vals := map[string]string{}
	for _, para := range consts_mappings {
		old_vari, new_vari := para.first, para.second
		if v, ok := vars_vals[old_vari]; ok {
			if v != new_vari {
				return false
			}
		} else {
			vars_vals[old_vari] = new_vari
		}
	}

	// Замена связанных переменных
	for vari, num := range new_vars {
		new_name := fmt.Sprintf("&%d", num)
		rs.substitute(vari, new_name, false)
	}
	// Замена констант
	for old_vari, new_vari := range vars_vals {
		rs.substitute(old_vari, new_vari, true)
	}

	return true
}

func (rs *ResolutionSearch) print_disjuncts() {
	log.Println("Дизъюнкты:")
	for _, d := range rs.disjuncts {
		log.Printf("  %s\n", d)
	}

	log.Println("--------------------")
}
