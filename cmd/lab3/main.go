package main

import (
	"expert_systems/pkg/models/and_or_tree"
	"log"
)

func main() {
	tr, err := and_or_tree.NewTree(
		and_or_tree.RuleFormat{
			Number:        101,
			InputsNumbers: []int{1, 2},
			ResultNumber:  11,
		},
		and_or_tree.RuleFormat{
			Number:        102,
			InputsNumbers: []int{3, 4},
			ResultNumber:  11,
		},
		and_or_tree.RuleFormat{
			Number:        103,
			InputsNumbers: []int{5, 6},
			ResultNumber:  12,
		},
		and_or_tree.RuleFormat{
			Number:        104,
			InputsNumbers: []int{7, 8},
			ResultNumber:  14,
		},
		and_or_tree.RuleFormat{
			Number:        105,
			InputsNumbers: []int{9, 10},
			ResultNumber:  21,
		},
		and_or_tree.RuleFormat{
			Number:        106,
			InputsNumbers: []int{11, 4, 12},
			ResultNumber:  19,
		},
		and_or_tree.RuleFormat{
			Number:        107,
			InputsNumbers: []int{12, 13, 14},
			ResultNumber:  20,
		},
		and_or_tree.RuleFormat{
			Number:        108,
			InputsNumbers: []int{15, 16},
			ResultNumber:  20,
		},
		and_or_tree.RuleFormat{
			Number:        109,
			InputsNumbers: []int{17, 18},
			ResultNumber:  22,
		},
		and_or_tree.RuleFormat{
			Number:        110,
			InputsNumbers: []int{19, 20},
			ResultNumber:  23,
		},
		and_or_tree.RuleFormat{
			Number:        111,
			InputsNumbers: []int{20, 21},
			ResultNumber:  23,
		},
		and_or_tree.RuleFormat{
			Number:        112,
			InputsNumbers: []int{21, 22},
			ResultNumber:  24,
		},
	)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(tr)

	// stck := stack.NewStack[*node.Node]()
	// log.Println("Поиск в глубину:")
	// algDFS := dfs.NewDeepSearch(gr, stck)
	// source := gr.Nodes[0]
	// target := gr.Nodes[8]
	// path, err := algDFS.FindTarget(source, target)

	// log.Println(path, err)

	// for _, v := range path {
	// 	log.Println(v.Number)
	// }

	// log.Println("Поиск в ширину:")
	// qu := queue.NewQueue[*node.Node]()
	// algBFS := bfs.NewWideSearch(gr, qu)
	// source = gr.Nodes[1]
	// target = gr.Nodes[8]
	// path, err = algBFS.FindTarget(source, target)

	// log.Println(path, err)
	// for _, v := range path {
	// 	log.Println(v.Number)
	// }
}
