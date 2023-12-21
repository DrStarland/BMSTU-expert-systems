package logic

import (
	"expert_systems/pkg/models/enums"
	"fmt"
)

type Rule struct {
	Id     int
	Inputs []Predicate
	Result Predicate
	Proved bool
}

func (rule Rule) String() string {
	inp := ""
	sep := " и "
	if len(rule.Inputs) > 1 {
		inp += rule.Inputs[0].String() + sep
	}
	for i := 0; i < len(rule.Inputs); i++ {
		inp += rule.Inputs[i].String() + sep
	}
	return fmt.Sprintf("№%d: если %sто %s", rule.Id, inp, rule.Result.String())
}

// ////////////////////////////////////////
type Variable struct {
	Name     string
	Negative bool
	Const    bool
	Status   enums.VarStatusEnum
}

func (t1 Variable) EqualTo(t2 Variable) bool {
	return (t1.Const && t2.Const) && (t1.Name == t2.Name)
}
func NewConst(name string) Variable    { return Variable{Name: name, Const: true} }
func NewVariable(name string) Variable { return Variable{Name: name, Const: false} }
func (t Variable) String() string      { return t.Name }

type Predicate struct {
	Name     string
	Args     []Variable
	Negative bool
	Proved   bool
}

func (p1 Predicate) EqualTo(p2 *Predicate) bool {
	if len(p1.Args) != len(p2.Args) {
		return false
	}

	return (p1.Name == p2.Name)
}

func (pred Predicate) String() string {
	res := ""
	if pred.Negative {
		res += "!"
	}

	res += pred.Name + "("
	first := true
	for _, v := range pred.Args {
		if !first {
			res += ", "
		}
		res += v.String()
		first = false
	}
	return res + ")"
}

type Disjunct struct {
	Predicates []*Predicate
}

func (d1 Disjunct) EqualTo(d2 Disjunct) bool {
	if len(d1.Predicates) != len(d2.Predicates) {
		return false
	}
	for i := range d1.Predicates {
		if !d1.Predicates[i].EqualTo(d2.Predicates[i]) {
			return false
		}
	}
	return true
}

func (d Disjunct) String() (res string) {
	if len(d.Predicates) == 0 {
		return "<empty>"
	}

	first := true
	for _, a := range d.Predicates {
		if !first {
			res += " V "
		}
		res += a.String()
		first = false
	}
	return res
}

func NewDisjunct(preds []*Predicate) Disjunct {
	if preds != nil {
		return Disjunct{Predicates: preds}
	}
	return Disjunct{Predicates: []*Predicate{}}
}

type Formula struct {
	Items []*Disjunct
}

// type Predicate_r struct {
// 	Name     string
// 	Negative bool
// 	Vars     []Var_r
// }

// func NewPredicate_r(p Predicate, disjunct_idx int, vars_dict map[string]Var_r) Predicate_r {
// 	self := Predicate_r{}

// 	self.Name = p.Name
// 	self.Negative = p.Negative

// 	if vars_dict == nil {
// 		vars_dict = make(map[string]Var_r)
// 	}

// 	self.Vars = make([]Var_r, 0)
// 	for _, arg := range p.Args {
// 		name := arg.Name
// 		if !unicode.IsUpper(rune(arg.Name[0])) {
// 			name = arg.Name + fmt.Sprintf("_%d", disjunct_idx)
// 		}
// 		if _, ok := vars_dict[name]; !ok {
// 			variable := NewVar_r(arg, name)
// 			vars_dict[name] = variable

// 		}

// 		self.Vars = append(self.Vars, vars_dict[name])
// 	}
// 	return self
// }

// func (pred Predicate_r) String() string {
// 	arguments := make([]string, 0, len(pred.Vars))
// 	for _, vr := range pred.Vars {
// 		arguments = append(arguments, vr.String())
// 	}
// 	x := fmt.Sprintf(`%s(%s)`, pred.Name, strings.Join(arguments, ", "))

// 	if !pred.Negative {
// 		return x
// 	} else {
// 		return string('¬') + x
// 	}
// }
