package formalaparsing

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/logic"
	"expert_systems/pkg/models/types"
	"log"
	"regexp"
	"strings"
)

func CleansedRunestrings(facts types.Runestring) []types.Runestring {
	facts_result := make([]types.Runestring, 0)
	t1 := facts.Split('\n')
	for _, str := range t1 {
		if len(str) > 0 {
			if _hm := types.Runestring(strings.TrimSpace(string(str))); len(_hm) > 0 {
				facts_result = append(facts_result, _hm)
			}
		}
	}
	return facts_result
}

// (Quantor | Op | Any | Predicate | Var)
func Full_parse_formula(f types.Runestring) logic.Term {
	f = types.Runestring(strings.ReplaceAll(string(f), " ", ""))
	res, _ := parse_formula_rec([]rune(f), nil)
	return res
}

func read_symbol(f types.Runestring) (enums.OperationTypeEnum, int) {
	s := f[0]
	optype, ok := enums.MapSymbolToOperationType[s]
	if !ok {
		log.Panic("incorrect symbol", s)
	}
	return optype, 1
}

var alphanumericCheck = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func read_var(f []rune) (interface{}, int) {
	s := f[0]
	if s == '(' {
		op, i := parse_formula_rec(f[1:], nil)
		return op, i + 2
	}

	// иначе - предикат или переменная
	if s == enums.SYMBOLS[enums.NOT] {
		res, i := read_var(f[1:])
		return logic.Op{
			Type: enums.NOT,
			Args: []any{res.(logic.Quantor)},
		}, i + 1
	}

	name := ""
	i := 0

	for i < len(f) && alphanumericCheck.MatchString(string(f[i])) && !enums.SYMBOLS.Include(f[i]) { // 3я проверка, ибо V
		name += string(f[i])
		i++
	}
	if i < len(f) && string(f[i]) == "(" {
		// это предикат
		i++
		args_str := ""

		for f[i] != ')' {
			args_str += string(f[i])
			i++
		}
		args := strings.Split(args_str, ",")

		args2 := make([]logic.Variable, 0)
		for _, arg_name := range args {
			args2 = append(args2, logic.Variable{Name: arg_name})
		}
		return logic.Predicate{
			Name:     name,
			Args:     args2,
			Negative: false,
		}, i + 1
	} else {
		return logic.Variable{
			Name:     name,
			Negative: false,
		}, i
	}
}

func parse_formula_rec(f []rune, arg1 interface{}) (logic.Term, int) {
	s := f[0]
	// 1-местные - операция, потом операнд
	if s == enums.SYMBOLS[enums.ALL] || s == enums.SYMBOLS[enums.EXISTS] {
		op, i := parse_formula_rec(f[2:], nil)
		variable := string(f[1])
		op_type := enums.MapSymbolToOperationType[s]
		return logic.Quantor{
			Type:      op_type,
			Variable:  variable,
			Operation: op.(logic.Op),
		}, i + 2
	}

	if s == enums.SYMBOLS[enums.NOT] {
		op, i := parse_formula_rec(f[1:], nil)
		return logic.Op{
			Type: enums.NOT,
			Args: []any{op.(logic.Quantor)},
		}, i + 1
	}

	// или переменная, или предикат

	// 2-местные
	idx1 := 0
	if arg1 == nil {
		arg1, idx1 = read_var(f)
	}

	if len(f) == idx1 || f[idx1] == ')' {
		// всё-таки одноместное...
		return arg1, idx1
	}

	op_type, idx2 := read_symbol(f[idx1:])
	arg2, idx3 := read_var(f[idx1+idx2:])

	idx_all := idx1 + idx2 + idx3
	op := logic.Op{
		Type: op_type,
		Args: []any{arg1, arg2},
	}

	if len(f) > idx_all && f[idx_all] != ')' && enums.SYMBOLS.Include(f[idx_all]) { // todo accurate
		op, idx_add := parse_formula_rec(f[idx_all:], op)
		return op, idx_all + idx_add
	} else {
		return op, idx_all
	}
}
