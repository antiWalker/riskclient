package core

// 字段类型
type FieldInterface interface {
	// 字段的名称
	Name() string
	// 字段支持的 操作符类型
	Operator() []AOT
}

/////////////////////////////////////////////////////////////////

// 获取符号的名称
func (ft FieldType) Name() string {
	return string(ft)
}

// 获取支持的操作符号
func (ft FieldType) Operator() []AOT {
	if l, exists := fieldSupportOperationMap[ft]; exists {
		return l
	} else {
		return []AOT{}
	}
}

type exprType string

const (
	exprReturnBool       exprType = "bool"
	exprReturnNumber     exprType = "number"
	exprReturnDuration   exprType = "duration"
	exprReturnDate       exprType = "date"
	exprReturnString     exprType = "string"
	exprReturnListNumber exprType = "string_of_number"
)
