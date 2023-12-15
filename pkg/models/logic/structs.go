package logic

import (
	"expert_systems/pkg/models/enums"
	"fmt"
	"strings"
	"unicode"
)

type Term interface {
	// Predicate_r | Rule_r
	// String() string
}

type Quantor struct {
	Type      enums.OperationTypeEnum
	Variable  string
	Operation Op
}

type Op struct {
	Type enums.OperationTypeEnum
	Args []any
}

type Variable struct {
	Name     string
	Negative bool
}

type Predicate struct {
	Name     string
	Args     []Variable
	Negative bool
}

type Var_r struct {
	Name   string
	Status enums.VarStatusEnum
	Value  string
}

func NewVar_r(v Variable, special_name string) Var_r {
	self := Var_r{}
	if special_name == "" {
		self.Name = v.Name
	} else {
		self.Name = special_name
	}
	if unicode.IsUpper(rune(self.Name[0])) {
		self.Status = enums.CONST
	} else {
		self.Status = enums.NOVAL
	}

	if unicode.IsUpper(rune(self.Name[0])) {
		self.Value = self.Name
	} else {
		self.Value = ""
	}
	return self
}

func (v Var_r) String() string {
	// if v.Value == "" || v.Status == enums.CONST {
	// 	return v.Name
	// }
	return fmt.Sprintf("%s=%s", v.Name, v.Value)
}

func (var_r Var_r) Resolution_eq(v Variable) bool {
	return var_r.Name == v.Name
}

type Rule_r struct {
	Id     int
	Inputs []Predicate_r
	Result Predicate_r
}

func (rule Rule_r) String() string {
	inp := ""
	for _, x := range rule.Inputs {
		inp += "и " + x.String() + ", "
	}
	return fmt.Sprintf("№%d: если %sто %s", rule.Id, inp, rule.Result.String())
}

type Predicate_r struct {
	Name     string
	Negative bool
	Vars     []Var_r
}

func NewPredicate_r(p Predicate, disjunct_idx int, vars_dict map[string]Var_r) Predicate_r {
	self := Predicate_r{}

	self.Name = p.Name
	self.Negative = p.Negative

	if vars_dict == nil {
		vars_dict = make(map[string]Var_r)
	}

	self.Vars = make([]Var_r, 0)
	for _, arg := range p.Args {
		name := arg.Name
		if !unicode.IsUpper(rune(arg.Name[0])) {
			name = arg.Name + fmt.Sprintf("_%d", disjunct_idx)
		}
		if _, ok := vars_dict[name]; !ok {
			variable := NewVar_r(arg, name)
			vars_dict[name] = variable

		}

		self.Vars = append(self.Vars, vars_dict[name])
	}
	return self
}

func (pred Predicate_r) String() string {
	arguments := make([]string, 0, len(pred.Vars))
	for _, vr := range pred.Vars {
		arguments = append(arguments, vr.String())
	}
	x := fmt.Sprintf(`%s(%s)`, pred.Name, strings.Join(arguments, ", "))

	if !pred.Negative {
		return x
	} else {
		return string('¬') + x
	}
}
