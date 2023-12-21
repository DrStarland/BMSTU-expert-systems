package DFS_logic

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/logic"
	"log"

	"github.com/emirpasic/gods/sets/hashset"
)

type Unifier struct {
	bindings               map[string]any // реально это срез строк
	substitutions          map[string]any
	variables              map[string]logic.Var_r
	transactions           []any
	forbiddenSubstitutions map[string]any
}

func NewUnifier(facts []logic.Predicate_r, rules []logic.Rule_r, target logic.Predicate_r) Unifier {
	un := Unifier{
		bindings:               make(map[string]any),
		substitutions:          make(map[string]any),
		variables:              make(map[string]logic.Var_r, len(facts)+2*len(rules)),
		transactions:           make([]any, 0),
		forbiddenSubstitutions: make(map[string]any),
	}

	for _, r := range rules {
		for _, a := range r.Inputs {
			for _, v := range a.Vars {
				un.variables[v.Name] = v
			}
			for _, v := range r.Result.Vars {
				un.variables[v.Name] = v
			}
		}
	}

	for _, f := range facts {
		for _, v := range f.Vars {
			un.variables[v.Name] = v
		}
	}
	for _, v := range target.Vars {
		un.variables[v.Name] = v
	}
	return un
}

func (un *Unifier) Substitute(var_name string, value string) {
	// if !custom_store {
	// 	store[var_name] = value
	// 	vari := un.variables[var_name]
	// 	vari.Value = value
	// 	vari.Status = enums.VAL

	// 	un.variables[var_name] = vari
	// } else {
	store := un.substitutions

	log.Println("SUBSTITUTIONS IN", var_name, value, store)

	if _, ok := store[var_name]; !ok {
		store[var_name] = []string{}
	}
	store[var_name] = append(store[var_name].([]string), value)
	// }

	if _, ok := un.substitutions[var_name]; ok {
		log.Println(un.variables[var_name])
		un.variables[var_name] = func(vari logic.Var_r) logic.Var_r {
			vari.Value = value
			vari.Status = enums.VAL
			return vari
		}(un.variables[var_name])
		log.Println(un.variables[var_name])

		for _, binded_var_name := range un.substitutions[var_name].([]string) {
			un.variables[binded_var_name] = func(vari logic.Var_r) logic.Var_r {
				vari.Value = value
				vari.Status = enums.VAL
				return vari
			}(un.variables[binded_var_name])
		}
	}

	log.Println("SUBSTITUTIONS OUT", un.variables)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func (un *Unifier) Bind(var_name1 string, var_name2 string) {
	if _, ok := un.bindings[var_name1]; !ok {
		un.bindings[var_name1] = make([]string, 0)
	}

	if _, ok := un.bindings[var_name2]; !ok {
		un.bindings[var_name2] = make([]string, 0)
	}

	un.bindings[var_name1] = append(un.bindings[var_name1].([]string), var_name2)
	un.bindings[var_name2] = append(un.bindings[var_name2].([]string), var_name1)
	// быстрый для написания способ удалить повторяющиеся значения
	// un.bindings[var_name1] = hashset.New(un.bindings[var_name1]).Values()
	// un.bindings[var_name2] = hashset.New(un.bindings[var_name2]).Values()

	un.bindings[var_name1] = removeDuplicateStr(un.bindings[var_name2].([]string))
	un.bindings[var_name1] = removeDuplicateStr(un.bindings[var_name2].([]string))

	log.Println("SET WORKS RESULT", un.bindings[var_name1], un.bindings[var_name2])

	// нельзя напрямую отредактировать поля структуры в карте, поэтому так
	binded_var := un.variables[var_name1]
	binded_var.Status = enums.BIND
	un.variables[var_name1] = binded_var

	binded_var = un.variables[var_name2]
	binded_var.Status = enums.BIND
	un.variables[var_name2] = binded_var
}

func (un *Unifier) applyChanges(substs, binds map[string]string) {
	log.Println("APPLY CHANGES", substs, binds)
	for var_name1, var_name2 := range binds {
		un.Bind(var_name1, var_name2)
	}
	for var_name, value := range substs {
		un.Substitute(var_name, value)
	}
}

// не реализовано
func (un *Unifier) cancelChanges(transactions []any) {
	if transactions == nil {
		return
	}
	for _, t := range transactions {
		un.cancelChange(t.([]any))
	}

}

func removeFromSlice(arr []string, tar string) []string {
	for i, st := range arr {
		if st == tar {
			if i != len(arr)-1 {
				return arr[:i]
			}
			return append(arr[:i], arr[i+1:]...)
		}
	}
	// append(slice[:index], slice[index+1:]...)
	return arr
}

func (un *Unifier) cancelChange(t []any) {
	substs, binds := t[0].(map[string]string), t[1].(map[string]string)

	for var_name1, var_name2 := range binds {
		un.bindings[var_name1] = removeFromSlice(un.bindings[var_name1].([]string), var_name2)
		un.bindings[var_name2] = removeFromSlice(un.bindings[var_name2].([]string), var_name1)
	}
	for var_name, _ := range substs {
		vari := un.variables[var_name]
		vari.Value = ""
		vari.Status = enums.NOVAL
		un.variables[var_name] = vari
		delete(un.substitutions, var_name)
	}
}

func (un *Unifier) unifyPredicate(a logic.Predicate_r, b logic.Predicate_r, check_only bool) (bool, []any) {
	// log.Println("Unification: ", a, b)
	canBeUnified := true
	binds := make(map[string]string)
	substs := make(map[string]string)
	if a.Name != b.Name {
		canBeUnified = false
	} else if len(a.Vars) != len(b.Vars) {
		canBeUnified = false
	} else {
		for i := range a.Vars {
			canBeUnified, bind, subst := un.unifyVeriable(a.Vars[i], b.Vars[i], true)
			// в одну переменную влить две подстановки нельзя: p(A, B) ~ p(x, x)
			if un.substitutionConficts(subst, un.substitutions) || un.checkSubstitutionForbidden(subst) {
				canBeUnified = false
			}
			if !canBeUnified {
				break
			}

			for key, value := range bind {
				binds[key] = value
			}
			for key, value := range subst {
				substs[key] = value
			}
		}
	}

	if check_only {
		return canBeUnified, nil
	}
	log.Println("ЧТО ЕСТЬ UNIFY PRED ", substs, binds)
	if canBeUnified {
		un.applyChanges(substs, binds)
	}

	return canBeUnified, []any{substs, binds}
}

func (un *Unifier) unifyVeriable(a, b logic.Var_r, check_only bool) (bool, map[string]string, map[string]string) {
	canBeUnified, binds, subst := true, make(map[string]string), make(map[string]string)

	switch {
	// обе переменные имеют значения: константа и константа, переменная и константа, переменная и переменная
	case a.Value != "" && b.Value != "":
		canBeUnified = a.Value == b.Value
	// только одна переменная имеет значение
	case a.Value != "" || b.Value != "":
		if a.Value != "" {
			a, b = b, a
		}
		subst[a.Name] = b.Value
	// обе переменные без значений
	default:
		binds[b.Name] = a.Name
		binds[a.Name] = b.Name
	}

	if check_only {
		return canBeUnified, binds, subst
	}

	un.applyChanges(subst, binds)
	return canBeUnified, binds, subst
}

func (un *Unifier) substitutionConficts(a map[string]string, b map[string]any) bool {
	overlap_keys := findKeysIntersection(a, b)
	log.Println(overlap_keys)
	// допустимо, когда подстановки направлены на одну переменную, но они одинаковы
	same := true
	for _, k := range overlap_keys {
		same = same && (a[k] == b[k])
		if !same {
			break
		}
	}
	log.Println("Иногда я возвращаю ложь", !same)
	return !same
}

func (un *Unifier) checkSubstitutionForbidden(subst map[string]string) bool {
	ok := true

	for var_name, val := range subst {
		if un.forbiddenSubstitutions[var_name] == nil {
			continue
		}

		forbidden := un.forbiddenSubstitutions[var_name]
		in_forbidden := false
		for _, x := range forbidden.([]string) {
			if x == val {
				in_forbidden = true
				break
			}
		}
		ok = ok && in_forbidden
	}
	log.Println("ЧТО ПРОИСХОДИТ: ", subst, ok)
	return !ok
}

func findKeysIntersection(a map[string]string, b map[string]any) []string {
	keys_a := getKeysAsAny(a)
	keys_b := getKeysAsAny(b)

	hs1 := hashset.New(keys_a...)
	hs2 := hashset.New(keys_b...)

	overlap_keys := hs1.Intersection(hs2).Values()
	result := make([]string, len(overlap_keys))
	for i, v := range overlap_keys {
		result[i] = v.(string)
	}

	log.Println("Overlapping", a, b, result)
	return result
}
