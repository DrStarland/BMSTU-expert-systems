package resolution

import (
	"expert_systems/pkg/models/logic"
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

func (rs *ResolutionSearch) Solve() {
	rs.print_disjuncts()
	decisionCanBeFound := true
	decisionFound := false

	for decisionCanBeFound && !decisionFound {
		decisionCanBeFound = false

	iteration:
		for i1, d1 := range rs.disjuncts {
			for i2, d2 := range rs.disjuncts {
				if d1 != d2 {
					decisionCanBeFound, decisionFound = rs.findOppositePair(rs.disjuncts[i1], rs.disjuncts[i2])
					if decisionCanBeFound {
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
