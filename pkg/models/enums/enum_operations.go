package enums

//go:generate go run golang.org/x/tools/cmd/stringer@v0.1.5 -type=OperationTypeEnum
type OperationTypeEnum int

// const SYMBOLS = "∀∃V&¬→~"

// OP_SYMBOL = lambda op: SYMBOLS[op.value]
// OP_SYMBOL_REV = lambda s: OpType(SYMBOLS.index(s))

type OperationSymbols []rune

func (base OperationSymbols) Include(s rune) bool {
	for _, r := range base {
		if r == s {
			return true
		}
	}
	return false
}

var SYMBOLS = OperationSymbols([]rune{'∀', '∃', 'V', '&', '¬', '→', '~'})

const (
	ALL OperationTypeEnum = iota
	EXISTS
	OR
	AND
	NOT
	IMPLY
	EQ
)

var MapSymbolToOperationType = func() map[rune]OperationTypeEnum {
	m := make(map[rune]OperationTypeEnum)
	for i := ALL; i <= EQ; i++ {
		m[SYMBOLS[i]] = i
	}
	return m
}()

var MapEnumStringToOperationType = func() map[string]OperationTypeEnum {
	m := make(map[string]OperationTypeEnum)
	for i := ALL; i <= EQ; i++ {
		m[i.String()] = i
	}
	return m
}()

// func (ee *EdgeLabelEnum) Scan(value interface{}) error {
// 	tempus := string(value.([]uint8))
// 	*ee = MapEnumStringToEntityStatus[tempus]
// 	return nil
// }
