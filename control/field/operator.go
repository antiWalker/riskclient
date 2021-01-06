package field

import "bigrisk/core"

/// 获取 fieldName 支持的 操作符号
func Operators(fieldName string) []string {
	ft := core.FieldType(fieldName)

	ops := ft.Operator()

	result := make([]string, len(ops))
	for idx, op := range ops {
		result[idx] = string(op)
	}

	return result
}
