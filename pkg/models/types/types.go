package types

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
