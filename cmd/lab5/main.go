package main

import (
	"expert_systems/pkg/algorithms/DFS_logic"
	formalaparsing "expert_systems/pkg/algorithms/formala-parsing"
	"expert_systems/pkg/models/types"
	"log"
)

var Runes = types.Runestring("ᚠ ᚢ ᚦ ᚫ ᚱ ᚲ ᚷ ᚹ ᚺ ᚾ ᛁ ᛃ ᛇ ᛈ ᛉ ᛋ ᛏ ᛒ ᛖ ᛗ ᛚ ᛝ ᛟ ᛞ ᚸ")

func main() {
	// ограничения: без вложенных предикатов, без бэктрекинга
	facts, rules, target := ex1()
	alg := DFS_logic.NewSearch(facts, rules, target)

	proved := alg.ProveTarget()
	log.Println(proved)
}

func ex1() ([]types.Runestring, []types.Runestring, types.Runestring) {
	facts := types.Runestring(
		`p1(W)
p6(M)
p7(N, M)
p8(N, A)`)
	rules := types.Runestring(
		`p1(x) & p2(y) & p3(x,y,z) & p4(z) → p5(x)
p6(x) & p7(N, x) → p3(W, x, N)
p6(x) → p2(x)
p8(x, A) → p4(x)
`)

	target := types.Runestring("p5(W)")
	facts_result := formalaparsing.CleansedRunestrings(facts)
	rules_result := formalaparsing.CleansedRunestrings(rules)
	return facts_result, rules_result, target
}

func ex2() ([]types.Runestring, []types.Runestring, types.Runestring) {
	facts := types.Runestring(
		`man(Adam)
man(Herasim)
man(Wallie)
man(Pup)
woman(Mumu)
woman(Eva)
child(Adam, Eva, Wallie)
child(Herasim, Mumu, Pup)`)
	rules := types.Runestring(
		`man(x) & child(x, y, z) → father(x, z)
man(x) & child(y, x, z) → father(x, z)
woman(x) & child(y, x, z) → mother(x, z)
woman(x) & child(x, y, z) → mother(x, z)
`)

	target := types.Runestring("mother(Eva, Wallie)")
	facts_result := formalaparsing.CleansedRunestrings(facts)
	rules_result := formalaparsing.CleansedRunestrings(rules)
	return facts_result, rules_result, target
}
