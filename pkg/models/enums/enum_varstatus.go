package enums

type VarStatusEnum int

const (
	CONST VarStatusEnum = iota
	NOVAL
	BIND
	VAL
)
