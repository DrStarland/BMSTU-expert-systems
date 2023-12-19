// Code generated by "stringer -type=OperationTypeEnum"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ALL-0]
	_ = x[EXISTS-1]
	_ = x[OR-2]
	_ = x[AND-3]
	_ = x[NOT-4]
	_ = x[IMPLY-5]
	_ = x[EQ-6]
}

const _OperationTypeEnum_name = "ALLEXISTSORANDNOTIMPLYEQ"

var _OperationTypeEnum_index = [...]uint8{0, 3, 9, 11, 14, 17, 22, 24}

func (i OperationTypeEnum) String() string {
	if i < 0 || i >= OperationTypeEnum(len(_OperationTypeEnum_index)-1) {
		return "OperationTypeEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OperationTypeEnum_name[_OperationTypeEnum_index[i]:_OperationTypeEnum_index[i+1]]
}