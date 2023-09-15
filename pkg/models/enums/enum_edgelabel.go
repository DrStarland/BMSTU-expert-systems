package enums

//go:generate go run golang.org/x/tools/cmd/stringer@v0.1.5 -type=EdgeLabelEnum
type EdgeLabelEnum int

const (
	Open EdgeLabelEnum = iota
	Closed
	Forbidden
)

var MapEnumStringToEdgeLabel = func() map[string]EdgeLabelEnum {
	m := make(map[string]EdgeLabelEnum)
	for i := Open; i <= Forbidden; i++ {
		m[i.String()] = i
	}
	return m
}()

// func (ee *EdgeLabelEnum) Scan(value interface{}) error {
// 	tempus := string(value.([]uint8))
// 	*ee = MapEnumStringToEntityStatus[tempus]
// 	return nil
// }
