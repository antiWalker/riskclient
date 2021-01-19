package core

import "time"

///// WARNING //////////////////////////////////////////////////////////
///// 在您修改代码之前请确保 完整的理解了代码
/////
///// 请确保所有 操作 符号的类型的 `值` 类型都不一样
/////
///// 没有特殊的原因需要加上这个限制
///// 虽然每个操作类型的 字段 和 值 可以唯一确认要做的符号
/////   [aka: dependent type | 上下文相关算法 ]
///// 但是我并不想实现[会增加整体的复杂度], 上下文无关算法 要比 上下文相关算法简单
/////
///// END OF WARNING ///////////////////////////////////////////////////

// AbstractOperatorType 抽象的操作符类型
type AbstractOperatorType string

//
type AOT = AbstractOperatorType

// 数学 运算 操作符 ----------------------------------------------
//  类型
type NumberOperatorType = AOT

// 数学 运算 操作符 ----------------------------------------------
// 只能有如下 几种
const (
	// decl not forget add to `AllNumberOperatorType`
	numberAdd      NumberOperatorType = "+"
	numberMinus    NumberOperatorType = "-"
	numberMultiply NumberOperatorType = "*"
	numberDivision NumberOperatorType = "/"
)

// 所有的 数学 运算操作符号
func AllNumberOperatorType() []NumberOperatorType {
	return []NumberOperatorType{
		numberAdd,
		numberMinus,
		numberMultiply,
		numberDivision,
	}
}

// 数学 运算 操作符号 结束 ------------------------------------

// 数学 比较 操作符 ------------------------------------------
type CompareOperatorType = AOT

const (
	// not forget add to `AllCompareOperatorType` type
	compareGreatThan  CompareOperatorType = ">"
	compareGreatEqual CompareOperatorType = ">="
	compareEqual      CompareOperatorType = "=="
	compareNotEqual   CompareOperatorType = "!="
	compareLessEqual  CompareOperatorType = "<="
	compareLessThan   CompareOperatorType = "<"
	compareBetween    CompareOperatorType = "between"  // [left, right) 左开右闭 区间
	compareInList     CompareOperatorType = "inIdList" // in list of id => eg:{1, 2, 3}
)

// 所有的 比较相关的 运算符
func AllCompareOperatorType() []CompareOperatorType {
	return []CompareOperatorType{
		compareGreatThan, compareGreatEqual, compareEqual, compareNotEqual,
		compareLessEqual, compareLessThan, compareBetween, compareInList,
	}
}

// 数学 比较 操作符 结束 --------------------------------------

// 逻辑操作符  -----------------------------------------------
// 类型
type LogicOperatorType = AOT

const (
	LogicAnd LogicOperatorType = "and"
	LogicOr  LogicOperatorType = "or"
	LogicNot LogicOperatorType = "not"
)

func AllLogicOperatorType() []LogicOperatorType {
	return []LogicOperatorType{LogicAnd, LogicOr, LogicNot}
}

// 逻辑操作符号 结束 ------------------------------------------

// 日期操作符 类型 --------------------------------------------
type DateOperatorType = AOT

const (
	DateBefore     DateOperatorType = "before"
	DateAfter      DateOperatorType = "after"
	DateInInterval DateOperatorType = "interval"
	DateMinus      DateOperatorType = "dateMinus"
)

func AllDateOperatorType() []DateOperatorType {
	return []DateOperatorType{DateBefore, DateAfter, DateInInterval, DateMinus}
}

// 日期 持续时间操作符号
// DateMinus op 的结果是此类型
// todo duration operation type
type DateDurationType time.Duration // we do use go system's type

////--------------------------------------------

// 日期操作符 结束 --------------------------------------------

// 字符串 操作类型 --------------------------------------------
type StringOperatorType = AOT

const (
	stringMinLen     StringOperatorType = "minLen"     // 最小长度
	stringMaxLen     StringOperatorType = "maxLen"     // 最大长度
	stringContain    StringOperatorType = "contain"    // 包含
	stringInWordList StringOperatorType = "inWordList" // 在一个 列表 中
	stringInDict     StringOperatorType = "inDict"     // 在一个 词库 中
)

func AllStringOperatorType() []StringOperatorType {
	return []StringOperatorType{stringContain, stringInWordList}
}

// 数据层 操作符类型 ------------------------------------------
type QueryOperatorType = AOT

const (
	queryMySQL QueryOperatorType = "mysql" // 查询 mysql
	queryES    QueryOperatorType = "es"    // 查询 es
	queryDict  QueryOperatorType = "dict"  // 查询 词库
	queryREDIS QueryOperatorType = "redis" // 查询 redis
)

func AllQueryOperatorType() []QueryOperatorType {
	return []QueryOperatorType{queryMySQL, queryES}
}

// 字符串操作类型 结束 -----------------------------------------

/// 操作符 返回的 类型
type (
	OpReturnType = exprType
	OpArgType    = exprType
)

/// 操作符 的信息
type OpInfo struct {
	ArgNumber  int          `json:"arg_number"`  // 这个操作符的参数
	ArgTypes   []OpArgType  `json:"arg_types"`   // len() == ArgNumber
	ReturnType OpReturnType `json:"return_type"` // return type
}

var allOpInfos = map[AOT]OpInfo{
	/// 数学 运算 操作符号
	numberAdd:      {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnNumber, exprReturnNumber}, ReturnType: exprReturnNumber}, // a + b
	numberMinus:    {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnNumber, exprReturnNumber}, ReturnType: exprReturnNumber}, // a - b
	numberMultiply: {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnNumber, exprReturnNumber}, ReturnType: exprReturnNumber}, // a * b
	numberDivision: {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnNumber, exprReturnNumber}, ReturnType: exprReturnNumber}, // a / b

	// 数值 比较 运算 操作符号
	// a >  b
	compareGreatThan: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */},
		ReturnType: exprReturnBool},

	// a >= b
	compareGreatEqual: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */},
		ReturnType: exprReturnBool},

	// a == b
	compareEqual: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */},
		ReturnType: exprReturnBool},

	// a <= b
	compareLessEqual: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */},
		ReturnType: exprReturnBool},
	// a <  b
	compareLessThan: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */},
		ReturnType: exprReturnBool},
	// a <= b < c {aka: b is in [a, c) }
	compareBetween: {ArgNumber: 3,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnNumber /* b */, exprReturnNumber /* c */},
		ReturnType: exprReturnBool},
	// a is member of list l {1, 2, 3}
	compareInList: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnNumber /* a */, exprReturnListNumber /* l */},
		ReturnType: exprReturnBool},

	// 逻辑 运算符号
	LogicAnd: {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnBool, exprReturnBool}, ReturnType: exprReturnBool}, // a && b
	LogicOr:  {ArgNumber: 2, ArgTypes: []OpArgType{exprReturnBool, exprReturnBool}, ReturnType: exprReturnBool}, // a || b
	LogicNot: {ArgNumber: 1, ArgTypes: []OpArgType{exprReturnBool /* */}, ReturnType: exprReturnBool},           // !a

	// 日期 运算符
	// event happened before a - b
	DateBefore: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnDate /* a */, exprReturnDuration /* b */},
		ReturnType: exprReturnBool},
	// event happened before a + b
	DateAfter: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnDate /* a */, exprReturnDuration /* b */},
		ReturnType: exprReturnBool},
	// a event happened in [b, c) [left open, right closed]
	DateInInterval: {ArgNumber: 3,
		ArgTypes:   []OpArgType{exprReturnDate /* a */, exprReturnDate /* b */, exprReturnDate /* c */},
		ReturnType: exprReturnBool},
	// duration after a event happened before b event happened
	DateMinus: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnDate /* a event */, exprReturnDate /* b event */},
		ReturnType: exprReturnDuration},

	// 字符串 运算符

	// a string have b sub-string
	stringContain: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnString /* a */, exprReturnString /* b */},
		ReturnType: exprReturnBool},

	// a string have any string in word list b
	stringInWordList: {ArgNumber: 2,
		ArgTypes:   []OpArgType{exprReturnString /* a */, exprReturnString /* b */},
		ReturnType: exprReturnBool},
}

// OperatorInfo return nil on error
func OperatorInfo(op AOT) *OpInfo {
	info, exists := allOpInfos[op]
	if exists {
		return &info // no use more memory
	}

	return nil
}
