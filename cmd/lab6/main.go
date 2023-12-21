package main

import (
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
	"fmt"
	"log"
)

var Runes = types.Runestring("ᚠ ᚢ ᚦ ᚫ ᚱ ᚲ ᚷ ᚹ ᚺ ᚾ ᛁ ᛃ ᛇ ᛈ ᛉ ᛋ ᛏ ᛒ ᛖ ᛗ ᛚ ᛝ ᛟ ᛞ ᚸ")

type Status int

const (
	NO_VALUE = iota
	LINKED
	HAS_VALUE
)

var hm logic.Variable

type Variable struct {
	name   string
	status Status
}

type Atom struct {
	name      string
	variables []Variable
	proved    bool
}

func (a Atom) Compare(b Atom) bool {
	if a.name != b.name {
		return false
	}

	if len(a.variables) != len(b.variables) {
		return false
	}

	return true
}

type Rule struct {
	id        int
	condition []Atom
	result    Atom
	proved    bool
}

type Unification struct {
	rules       []Rule
	closedRules mutableSetOfInt  // int
	closedAtoms mutableSetOfAtom // atom
}

type mutableSetOfInt struct {
	list []int
	_map map[int]bool
}
type mutableSetOfAtom struct {
	list []Atom
	_map map[*Atom]bool
}

func NewMutableSetOfInt() mutableSetOfInt {
	return mutableSetOfInt{
		list: make([]int, 0),
		_map: map[int]bool{},
	}
}

func (s *mutableSetOfInt) add(x int) {
	if _, ok := s._map[x]; !ok {
		s.list = append(s.list, x)
		s._map[x] = true
	}
}

func NewMutableSetOfAtom() mutableSetOfAtom {
	return mutableSetOfAtom{
		list: make([]Atom, 0),
		_map: map[*Atom]bool{},
	}
}

func NewUnification() Unification {
	return Unification{
		rules:       make([]Rule, 0),
		closedRules: NewMutableSetOfInt(),
		closedAtoms: NewMutableSetOfAtom(),
	}
}

func (s *mutableSetOfAtom) Add(elems ...Atom) {
	s.list = append(s.list, elems...)
	for _, el := range elems {
		s._map[&el] = true
	}
}

func (s *mutableSetOfAtom) append(x Atom) {
	if _, ok := s._map[&x]; !ok {
		s.list = append(s.list, x)
	}
}

func (s mutableSetOfAtom) Size() int {
	return len(s.list)
}

func (s mutableSetOfInt) Contains(target int) bool {
	_, ok := s._map[target]
	return ok
}

func (un *Unification) prove(facts []Atom, targetRule int) bool {
	un.closedAtoms.Add(facts...)

	var closedAtomsChanged = true
	for closedAtomsChanged && !un.closedRules.Contains(targetRule) {
		oldSize := un.closedAtoms.Size()
		un.research()
		closedAtomsChanged = un.closedAtoms.Size() != oldSize
	}

	log.Println("___________________________________________")
	log.Println(len(un.closedRules.list), un.closedRules.list)
	log.Println("Facts", un.closedAtoms.list)

	return un.closedRules.Contains(targetRule)
}

func (un *Unification) research() {
	for k, rule := range un.rules {
		if !rule.proved {
			log.Printf("*** Рассматриваем правило %v\n", rule)
			for i, atom := range rule.condition {
				if !atom.proved {
					for _, it := range un.closedAtoms.list {
						if atom.name == it.name && len(atom.variables) == len(it.variables) {
							log.Printf("Рассматриваем атом $atom правила %d\n", rule.id)
							un.rules[k].condition[i].variables = it.variables
							un.rules[k].condition[i].proved = true
							log.Printf("Доказан атом %v\n", atom)
							break
						}
					}
				} else {
					log.Printf("Атом %v правила %v уже доказан\n", atom, rule.id)
				}
			}

			all := true
			for _, hm := range rule.condition {
				if !hm.proved {
					all = false
					break
				}
			}
			if all {
				log.Printf("Правило %d доказано!", rule.id)
				for j, resVar := range rule.result.variables {
					// resVar.name = rule.condition
					// 	.flatMap { it.variables } // берем вообще все переменные условия правила
					// 	.first { it.name[0].lowercase() == resVar.name[0].lowercase() } // ищем наименование переменной
					// 	.name
					all_vars := make([]Variable, 0)
					for _, em := range rule.condition {
						all_vars = append(all_vars, em.variables...)
					}
					new_name := ""
					for _, it := range all_vars {
						if it.name[0] == resVar.name[0] {
							new_name = it.name
							break
						}
					}
					un.rules[k].result.variables[j].name = new_name
				}
				un.rules[k].proved = true
				log.Printf("Добавляем атом %v", rule.result)
				un.closedRules.add(rule.id)
				un.closedAtoms.append(rule.result)
			}
		}
	}
	log.Println()
	log.Println("В результате итерации список закрытых фактов: ", un.closedAtoms.list)
	log.Println("=====================")
}

func main() {
	a := NewUnification()
	a.rules = append(a.rules,
		Rule{
			id: 0,
			condition: []Atom{
				Atom{
					name:      "не умеет летать",
					variables: []Variable{Variable{name: "x"}},
				},
				Atom{
					name:      "ловит рыбу под водой",
					variables: []Variable{Variable{name: "x"}},
				},
			},
			result: Atom{
				name:      "умеет плавать",
				variables: []Variable{Variable{name: "x"}},
			},
		},
		Rule{
			id: 1,
			condition: []Atom{
				Atom{
					name:      "является пингвином",
					variables: []Variable{Variable{name: "x"}},
				},
			},
			result: Atom{
				name:      "ловит рыбу под водой",
				variables: []Variable{Variable{name: "x"}},
			},
		},
		Rule{
			id: 2,
			condition: []Atom{
				Atom{
					name:      "является пингвином",
					variables: []Variable{Variable{name: "x"}},
				},
			},
			result: Atom{
				name:      "является птицей",
				variables: []Variable{Variable{name: "x"}},
			},
		},
		Rule{
			id: 3,
			condition: []Atom{
				Atom{
					name:      "является пингвином",
					variables: []Variable{Variable{name: "x"}},
				},
			},
			result: Atom{
				name:      "не умеет летать",
				variables: []Variable{Variable{name: "x"}},
			},
		},
	)

	res := a.prove(
		[]Atom{ // facts =
			Atom{
				name:      "является пингвином",
				variables: []Variable{Variable{name: "x_GERTRUDA", status: HAS_VALUE}},
				proved:    true,
			},
			//            Atom{
			//                name = "не умеет летать",
			//                variables = []Variable{Variable{"X_GERTRUDA", Variable.Status.HAS_VALUE}}
			//            }
		},
		0, // targetRule
	)

	log.Println(res)

	aa := logic.Predicate{Name: "Kek", Args: []logic.Variable{logic.Variable{Name: "LOL", Const: false}}}
	bb := logic.Predicate{Name: "Kek", Args: []logic.Variable{logic.Variable{Name: "Cheburek", Const: false}}}
	cc := logic.Predicate{Name: "Kek", Args: []logic.Variable{logic.Variable{Name: "Shrek", Const: true}}}
	kb := []logic.Predicate{aa, bb, cc}
	log.Println(Unify(kb, &aa, &bb), "LB", kb, "AA", aa, "BB", bb)
	log.Println(Unify(kb, &bb, &cc), "LB", kb, "AA", aa, "BB", bb, "CC", cc)
}

func Unify(kb []logic.Predicate, a1, a2 *logic.Predicate) bool {
	type pair struct {
		first  string
		second string
	}

	if a1.Name != a2.Name {
		log.Panic(a1.Name, a2.Name)
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
			old_v, _ := para.first, para.second
			if old_v == vari {
				consts_mappings[i].first = fmt.Sprintf("@%d", num)
			}
		}
	}

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
		substitute(kb, vari, new_name, false)
	}
	// Замена констант
	for old_v, new_v := range vars_vals {
		substitute(kb, old_v, new_v, true)
	}

	return true
}

func substitute(kb []logic.Predicate, old_v, new_v string, make_const bool) {
	_const := ""
	if make_const {
		_const = " const"
	}
	log.Println("  Замена: " + old_v + " -> " + new_v + _const)

	for j := range kb {
		for k, vari := range kb[j].Args {
			if vari.Name == old_v {
				vari.Name, vari.Const = new_v, make_const
				kb[j].Args[k] = vari
			}
		}
	}
}

// func main() {
// 	// ограничения: без вложенных предикатов, без бэктрекинга
// 	facts, rules, target := ex2()
// 	alg := DFS_logic.NewSearch(facts, rules, target)

// 	proved := alg.ProveTarget()
// 	log.Println(proved)
// }

// func ex1() ([]types.Runestring, []types.Runestring, types.Runestring) {
// 	facts := types.Runestring(
// 		`p1(W)
// p6(M)
// p7(N, M)
// p8(N, A)`)
// 	rules := types.Runestring(
// 		`p1(x) & p2(y) & p3(x,y,z) & p4(z) → p5(x)
// p6(x) & p7(N, x) → p3(W, x, N)
// p6(x) → p2(x)
// p8(x, A) → p4(x)
// `)

// 	target := types.Runestring("p5(W)")
// 	facts_result := formalaparsing.CleansedRunestrings(facts)
// 	rules_result := formalaparsing.CleansedRunestrings(rules)
// 	return facts_result, rules_result, target
// }

// func ex2() ([]types.Runestring, []types.Runestring, types.Runestring) {
// 	facts := types.Runestring(
// 		`man(Adam)
// man(Herasim)
// man(Wallie)
// man(Pup)
// woman(Mumu)
// woman(Eva)
// child(Adam, Eva, Wallie)
// child(Herasim, Mumu, Pup)`)
// 	rules := types.Runestring(
// 		`man(x) & child(x, y, z) → father(x, z)
// man(x) & child(y, x, z) → father(x, z)
// woman(x) & child(y, x, z) → mother(x, z)
// woman(x) & child(x, y, z) → mother(x, z)
// `)

// 	target := types.Runestring("father(Herasim, Wallie)")
// 	facts_result := formalaparsing.CleansedRunestrings(facts)
// 	rules_result := formalaparsing.CleansedRunestrings(rules)
// 	return facts_result, rules_result, target
// }
