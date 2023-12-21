package types

import "strings"

type Runestring []rune

func (runeword Runestring) String() string {
	return string(runeword)
}

func (runeword Runestring) IndexOf(target rune) int {
	for i, runa := range runeword {
		if runa == target {
			return i
		}
	}
	return -1
}

func (runeword Runestring) Split(sep rune) []Runestring {
	result := make([]Runestring, 0)

	for len(runeword) > 0 {
		idx := runeword.IndexOf(sep)
		if idx < 0 {
			result = append(result, runeword)
			break
		}
		result = append(result, runeword[:idx])
		runeword = runeword[idx+1:]
	}

	return result
}

func CleansedRunestrings(facts Runestring) []Runestring {
	facts_result := make([]Runestring, 0)
	t1 := facts.Split('\n')
	for _, str := range t1 {
		if len(str) > 0 {
			if _hm := Runestring(strings.TrimSpace(string(str))); len(_hm) > 0 {
				facts_result = append(facts_result, _hm)
			}
		}
	}
	return facts_result
}
