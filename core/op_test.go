package core

import "testing"

func TestDemo(t *testing.T) {
	t.Log("this is a demo test")
}

/// 检测是否有重复定义的运算符
func TestIfHaveSameOperator(t *testing.T) {

	// get all ops
	number := AllNumberOperatorType()
	compare := AllCompareOperatorType()
	logic := AllLogicOperatorType()
	date := AllDateOperatorType()
	str := AllStringOperatorType()

	//////////////////////////////////////////////////////////////////////////
	ops := make(map[AbstractOperatorType]interface{}) // fuck off no `set` type

	// do check if empty
	addToOps := func(newOps []AbstractOperatorType) {
		for _, op := range newOps {
			if _, exists := ops[op]; exists {
				t.Error("try to redefine ", op, " !", " FATAL ERROR ! [read readme.md before add new operator] ")
			} else {
				ops[op] = nil
			}
		}
	}

	addToOps(number)
	addToOps(compare)
	addToOps(logic)
	addToOps(date)
	addToOps(str)
	//////////////////////////////////////////////////////////////////////////
}
