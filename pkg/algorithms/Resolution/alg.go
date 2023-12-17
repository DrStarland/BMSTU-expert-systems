package resolution

import (
	"expert_systems/pkg/models/logic"
	"log"
)

type ResolutionSolver struct {
	disjuncts    []*logic.Disjunct
	final_result bool
	iter_changed bool
}

// "конструктор" алгоритма поиска
func NewSearch(formulas []logic.Formula, neg_target logic.Formula) ResolutionSolver {
	rs := ResolutionSolver{
		disjuncts: make([]*logic.Disjunct, 0, len(formulas)+len(neg_target.Items)),
	}

	for _, f := range formulas {
		rs.disjuncts = append(rs.disjuncts, f.Items...)
	}
	rs.disjuncts = append(rs.disjuncts, neg_target.Items...)
	return rs
}

// void solve() {
// 	while (iter_changed) {
// 		for (Disjunct &d1 : disjuncts) {
// 			for (Disjunct &d2 : disjuncts) {
// 				if (&d1 == &d2)
// 					continue;

// 				bool match = check_disjuncts(d1, d2);
// 				if (match) {
// 					iter_changed = true;
// 					break;
// 				}
// 			}
// 		}
// 	}

// }

func (rs *ResolutionSolver) Solve() {
	rs.print_disjuncts()
	rs.iter_changed = true
	for rs.iter_changed {
		rs.iter_changed = false

		for i1, d1 := range rs.disjuncts {
			for i2, d2 := range rs.disjuncts {
				if d1 == d2 {
					log.Println("ХОТЬ РАЗ ЭТО РАБОТАЛО?")
					continue
				}

				match := rs.check_disjuncts(rs.disjuncts[i1], rs.disjuncts[i2])
				if match {
					rs.iter_changed = true
					break
				}
			}
			if rs.iter_changed {
				break
			}
		}

		rs.print_disjuncts()
		if rs.final_result {
			break
		}

	}

	if rs.final_result {
		log.Println("Доказано")
	} else {
		log.Println("He доказано")
	}
}
