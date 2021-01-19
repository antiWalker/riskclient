package core

import (
	"bigrisk/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 节点类型
type nodeType string

const (
	integerNodeType  nodeType = "integer"  // 整数类型
	numberNodeType   nodeType = "number"   // 实数类型
	stringNodeType   nodeType = "string"   // 字符串节点
	dateNodeType     nodeType = "date"     // 日期类型节点
	durationNodeType nodeType = "duration" // 持续时间节点
	boolNodeType     nodeType = "bool"     // 布尔节点 [用于执行时记录子表达式的结果]
	fieldNodeType    nodeType = "field"    // 字段类型 [读取 InParams 指定的字段]
	selectNodeType   nodeType = "select"   // 数据库查询目标节点
	whereNodeType    nodeType = "where"    // 数据库查询条件节点

	operatorNodeType nodeType = "operator" // 运算操作节点
	logicNodeType    nodeType = "logic"    // 逻辑操作节点
	queryNodeType    nodeType = "query"    // 数据库查询节点
)

/// 简单的值节点
/// DO NOT have children
/// ```json
/// {
///     "type": "plainNumber",
///     "value": 2.3,           // ....
/// }
/// ```
type simpleNode struct {
	Type  nodeType    `json:"type"`
	Value interface{} `json:"value"`
}

func (s *simpleNode) String() string {
	str := fmt.Sprint(s.Value)

	return "type: " + string(s.Type) + "\n" + "value: " + str
}

/// ```json
/// {
///     "type": "logic",
///     "operator": "and",
///     "children": [
///     //... list of expr
///     ]
/// }
/// ```
///
/// ```json
/// {
///     "type": "numberOperator",
///     "operator": "+",
///     "children": [
///     //... list of expr
///     ]
/// }
/// ```
type complexNode struct {
	Type  nodeType             `json:"type"`
	Value AbstractOperatorType `json:"value"` // 操作符号的值
	/// type can be:
	/// * nil          there is sth error
	/// * int64        ok for direct-integer-value
	/// * float64      ok for direct-number-value
	/// * string       ok for direct-string-value
	/// * simpleNode   ok for direct-value
	/// * complexNode  ok for sub-expr
	Children []interface{} `json:"children"`
}

func (c *complexNode) String() string {

	str := "type : " + string(c.Type) + "\n" + "value: " + string(c.Value)

	l := make([]string, len(c.Children))
	for idx, ch := range c.Children {
		switch ch.(type) {
		case *simpleNode:
			l[idx] = ch.(*simpleNode).String()
		case *complexNode:
			l[idx] = ch.(*complexNode).String()
		default:
			l[idx] = ""
		}
	}

	str = str + "\nchildren:"
	for _, item := range l {
		lines := strings.Split(item, "\n")
		for _, line := range lines {
			str = str + "\n  " + line
		}
	}
	str = str + "\n"

	return str

}

// 规则类型
type rawRuleType struct {
	Sign      string                 `json:"sign"`
	Condition map[string]interface{} `json:"condition"`
	Match     map[string]interface{} `json:"match"`
	Exception map[string]interface{} `json:"exception"`
}

//
type cookedRuleType struct {
	Sign      string
	Condition *complexNode // nil if not exists
	Match     *complexNode // nil if not exists
	Exception *complexNode // nil if not exists
}

/// helper function
func deRefInterface(value reflect.Value) reflect.Value {
	for {
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		} else {
			break
		}
	}
	return value
}

/// return value is
/// * nil              on failed
/// * string           on nt = stringNodeType
/// * time.Duration    on nt = durationNodeType
/// * time.Time        on nt = dateNodeType
/// * int64            on nt = integerNodeType
/// * float64          on nt = numberNodeType
func parseSimpleNodeValue(nt nodeType, value interface{}) interface{} {

	switch nt {
	case durationNodeType:
		if s, ok := value.(string); ok {
			if d, err := time.ParseDuration(s); err == nil {
				return d
			}
		}
		return nil

	case stringNodeType:
		if s, ok := value.(string); ok {
			return s
		}
		return nil

	case fieldNodeType:
		if s, ok := value.(string); ok {
			return s
		}
		return nil

	case dateNodeType:
		if s, ok := value.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
			}
		}
		return nil
	case integerNodeType:
		return int64(value.(float64))
	case numberNodeType:
		return value.(float64)
	case whereNodeType:
		if s, ok := value.(string); ok {
			return s
		}
		return nil
	case selectNodeType:
		if s, ok := value.(string); ok {
			return s
		}
		return nil
	default:
		return nil
	}
}

/// nt must be some simple node
/// aka must in following list:
/// * durationNodeType:
/// * dateNodeType
/// * integerNodeType
/// * numberNodeType
/// * stringNodeType
/// * fieldNodeType
func constructSimpleNode(nt nodeType, value interface{}) (*simpleNode, error) {
	simple := new(simpleNode)

	simple.Type = nt

	v := parseSimpleNodeValue(nt, value)

	if v == nil {
		return nil, errors.New("construct simple node failed: " + string(nt))
	}

	simple.Value = v

	return simple, nil
}

/// nt must be some complex node
func constructComplexNode(nt nodeType, value interface{}, children interface{}) (*complexNode, error) {
	node := new(complexNode)

	node.Type = nt

	if _, ok := value.(string); !ok {
		return nil, errors.New("`value` field must be string type")
	}

	node.Value = AbstractOperatorType(value.(string))

	l := reflect.ValueOf(children)

	switch l.Kind() {
	case reflect.Array, reflect.Slice:
		length := l.Len()
		node.Children = make([]interface{}, length)

		for idx := 0; idx < length; idx++ {
			v := deRefInterface(l.Index(idx))

			switch v.Kind() {
			case reflect.Map:
				subExpr, err := tryConstructUnknownNode(v.Interface())
				if err == nil {
					node.Children[idx] = subExpr
				} else {
					return nil, errors.New("parse sub-expr failed:" + err.Error())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				node.Children[idx] = v.Interface()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				node.Children[idx] = v.Interface()
			}
		}

		return node, nil
	default:
		return nil, errors.New("`children` field must be array")
	}
}

func tryConstructUnknownNode(node interface{}) (interface{}, error) {
	nodeValue := reflect.ValueOf(node)

	nodeValue = deRefInterface(nodeValue)

	if nodeValue.Kind() != reflect.Map {
		return nil, errors.New("node is not map type")
	}

	// node have `map` type
	nodeKeys := nodeValue.MapKeys()

	nodeKVMap := map[string]reflect.Value{}
	for _, key := range nodeKeys {
		switch key.Kind() {
		case reflect.String:
			s := key.Interface().(string)
			nodeKVMap[s] = key
		default: // according to  json format std: key are MUST BE string
			// hence there is no chance to execute to here
			return nil, errors.New("node's key must be string type")
		}
	}
	////////////////////////////////////////////////////////////////////////////////

	getValue := func(key string) interface{} {
		keyValue, exists := nodeKVMap[key]
		if !exists {
			return nil
		}

		valueValue := nodeValue.MapIndex(keyValue)

		valueValue = deRefInterface(valueValue)

		return valueValue.Interface()
	}

	ntV := getValue("type")
	if ntV == nil {
		return nil, errors.New("`type` field is not exist or empty")
	}
	nt, ok := ntV.(string)
	if !ok {
		return nil, errors.New("`type` field is not string")
	}

	value := getValue("value")
	if value == nil {
		return nil, errors.New("`value` field is not exist or empty")
	}

	/// construct complex node type helper
	ch := func(nt nodeType) (interface{}, error) {
		childrenKeyValue := nodeKVMap["children"]

		if childrenKeyValue.IsValid() == false {
			return nil, errors.New("`children` is no valid:" + childrenKeyValue.String())
		}

		children := nodeValue.MapIndex(childrenKeyValue)

		children = deRefInterface(children)

		switch children.Kind() {
		case reflect.Array, reflect.Slice:
			return constructComplexNode(nt, value, children.Interface())
		default:
			return nil, errors.New("`children` have wrong type")
		}
	}

	switch nodeType(nt) {
	/// simple node
	case stringNodeType:
		return constructSimpleNode(stringNodeType, value)

	case dateNodeType:
		return constructSimpleNode(dateNodeType, value)

	case durationNodeType:
		return constructSimpleNode(durationNodeType, value)

	case integerNodeType:
		return constructSimpleNode(integerNodeType, value)

	case numberNodeType:
		return constructSimpleNode(numberNodeType, value)

	case fieldNodeType:
		return constructSimpleNode(fieldNodeType, value)

	case selectNodeType:
		return constructSimpleNode(selectNodeType, value)

	case whereNodeType:
		return constructSimpleNode(whereNodeType, value)
	///////////////////////////////////////////////////////////////////////

	/// complex node
	case operatorNodeType:
		return ch(operatorNodeType)

	case logicNodeType:
		return ch(logicNodeType)

	case queryNodeType:
		return ch(queryNodeType)
	///////////////////////////////////////////////////////////////////////

	default:
		return nil, errors.New("unknown field type:" + nt)
	}
}

/// return value is
/// * (*simpleNode, nil)
/// * (*complexNode, nil)
/// * (nil, error)
func tryConstructNode(m map[string]interface{}) (interface{}, error) {
	if _, ok := m["type"]; ok {
		node, err := tryConstructUnknownNode(m)

		return node, err
	}

	return nil, errors.New("node is empty")
}

func constructNodeFromString(ctx context.Context, ruleStr []byte) (*cookedRuleType, error) {
	var TraceId string
	if v := ctx.Value("TraceId"); v != nil {
		TraceId = strconv.Itoa(v.(int))
	}
	compileStart := time.Now().UnixNano()

	var raw rawRuleType

	cooked := new(cookedRuleType)

	// json decode failed
	if err := json.Unmarshal(ruleStr, &raw); err != nil {
		return nil, err
	}
	if raw.Sign == "" {
		return nil, errors.New("The Sign of Rule is Empty!\n ")
	}

	doC := func(m map[string]interface{}) (*complexNode, error) {
		// construct match
		match, err := tryConstructNode(m)
		if err != nil {
			return nil, err
		}

		cast, ok := match.(*complexNode)
		if !ok {
			return nil, err
		}

		return cast, nil
	}

	if raw.Condition != nil {
		n, err := doC(raw.Condition)
		if err != nil {
			return nil, err
		}
		cooked.Condition = n
	}

	if raw.Match != nil {
		n, err := doC(raw.Match)
		if err != nil {
			return nil, err
		}
		cooked.Match = n
	}

	if raw.Exception != nil {
		n, err := doC(raw.Exception)
		if err != nil {
			return nil, err
		}
		cooked.Exception = n
	}

	cooked.Sign = raw.Sign

	//compileElapsed := time.Since(compileStart)

	//log.Debug("DetectHandler Compile Cost Time: ", (time.Now().UnixNano()-compileStart)/1000,"costTime")
	common.InfoLogger.Infof(" TraceId : %v , Wrapper Cost Time: %v", TraceId, (time.Now().UnixNano()-compileStart)/1000)

	return cooked, nil
}
