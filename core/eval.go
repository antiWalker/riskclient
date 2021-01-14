package core

import (
	"bigrisk/models"
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
)

func (c *complexNode) Execute(ctx context.Context, runStack *Stack, params map[string]interface{}, reason *[]string) error {
	for _, ch := range c.Children {

		switch ch.(type) {
		case *simpleNode:
			// If this node is simple one, push to stack
			runStack.Push(ch.(*simpleNode))
		case *complexNode:
			ok := ch.(*complexNode).Execute(ctx, runStack, params, reason)

			if ok != nil {
				return errors.New("riskEngine: executeResult Error\n" + ok.Error())
			}

			middleResult, ok1 := ExecuteComplexNode(ctx, ch.(*complexNode), runStack, params, reason)

			if ok1 != nil {
				return errors.New("riskEngine: executeComplexNode Error\n" + ok1.Error())
			}

			runStack.Push(middleResult)
		default:
			return errors.New("riskEngine: not recogonized node type")
		}
	}

	return nil
}

// ExecuteComplexNode execute one complex node
func ExecuteComplexNode(ctx context.Context, c *complexNode, runStack *Stack, params map[string]interface{}, reason *[]string) (interface{}, error) {
	switch c.Type {
	case operatorNodeType:
		return ExecuteOperatorNode(ctx, c, runStack, params, reason)
	case queryNodeType:
		return ExecuteQueryNode(ctx, c, runStack, params, reason)
	case logicNodeType:
		return ExecuteLoginNode(ctx, c, runStack, params)
	default:
		return nil, errors.New("riskEngine: Execute op not support" + string(c.Type))
	}
}

// ExecuteQueryNode execute one query type node
func ExecuteQueryNode(ctx context.Context, c *complexNode, runStack *Stack, params map[string]interface{}, reason *[]string) (interface{}, error) {
	var TraceId string
	if v := ctx.Value("TraceId"); v != nil {
		TraceId = v.(string)
	}
	switch c.Value {
	case queryMySQL:

		var jobs []models.Job
		var wh []models.Where

		var tableStr string
		var columnStr string

		resultNodeType := integerNodeType

		//log.Info("runStack", runStack,"runStack")
		log.Info("runStack", &TraceContext{
			TraceId:  TraceId,
			RunStack: runStack,
		})
		for _, tmpNode := range *runStack {
			if tmpNode.(*simpleNode).Type == selectNodeType || tmpNode.(*simpleNode).Type == whereNodeType {

				whereStr := tmpNode.(*simpleNode).Value.(string)
				switch tmpNode.(*simpleNode).Type {
				case whereNodeType:

					whereArr := strings.Split(whereStr, "|")

					columnStr := whereArr[0]
					opStr := whereArr[1]
					valueStr := whereArr[2]

					//log.Debug("columnStr", whereArr,"columnStr")
					//log.Debug("opStr", whereArr,"opStr")
					//log.Debug("valueStr", whereArr,"valueStr")

					if strings.HasPrefix(valueStr, "$") {
						valueMap := strings.Split(valueStr, "$")
						//	valueStr := reflect.ValueOf(params).Elem().FieldByName(valueMap[1])
						key := valueMap[1]

						valueStr, ok := params[key]
						if ok == false {
							return nil, errors.New("riskEngine: valueStr is not valid\n ")
						}
						typ := reflect.TypeOf(valueStr)
						switch typ.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							valueReal := valueStr.(int64)
							wh = append(wh, models.Where{columnStr, opStr, strconv.FormatInt(valueReal, 10)})

						case reflect.Float32, reflect.Float64:
							valueReal := valueStr.(float64)
							wh = append(wh, models.Where{columnStr, opStr, strconv.FormatFloat(valueReal, 'f', -1, 64)})

						case reflect.String:
							if valueStr.(string) != "" {
								wh = append(wh, models.Where{columnStr, opStr, valueStr.(string)})
							}
						default:
							return nil, errors.New("riskEngine: valueStr's Kind is Not Supported\n")
						}
					} else if strings.HasPrefix(valueStr, "T") {
						valueMap := strings.Split(valueStr, "T")
						hourBefore := valueMap[1]

						var baseTime time.Time

						if BaseTime == "real_time" {
							baseTime = time.Now()
						} else {
							key := changeField(columnStr)

							if key == "" {
								return nil, errors.New("riskEngine: " + columnStr + " is empty \n")
							}

							timeStr, ok := params[key]
							if ok == false {
								return nil, errors.New("riskEngine: timeStr is not valid\n ")
							}

							var err error
							baseTime, err = time.Parse("2006-01-02 15:04:05", timeStr.(string))

							if err != nil {
								return nil, errors.New("riskEngine: " + columnStr + " is not a valid DateStr \n" + err.Error())
							}
						}

						timeBegin := baseTime
						timeEnd := baseTime

						switch hourBefore {
						case "24":
							timeBegin = baseTime.Add(time.Hour * -24)
						case "168":
							timeBegin = baseTime.Add(time.Hour * -24 * 7)
						case "360":
							timeBegin = baseTime.Add(time.Hour * -24 * 15)
						case "720":
							timeBegin = baseTime.Add(time.Hour * -24 * 30)
						default:
							timeBegin = baseTime.Add(time.Hour * -24 * 7)
						}

						wh = append(wh, models.Where{columnStr, opStr, strconv.FormatInt(timeBegin.Unix(), 10)})
						wh = append(wh, models.Where{columnStr, "lte", strconv.FormatInt(timeEnd.Unix(), 10)})

					} else {
						if valueStr != "" {
							wh = append(wh, models.Where{columnStr, opStr, valueStr})
						}
					}

				case selectNodeType:
					selectStr := tmpNode.(*simpleNode).Value.(string)

					selectArr := strings.Split(selectStr, "|")

					tableStr = selectArr[0]
					columnStr = selectArr[1]

					if strings.Contains(columnStr, "sum") {
						resultNodeType = numberNodeType
					}

				default:
					return nil, errors.New("riskEngine: unexpect type in query execute\n ")
				}

				var err error
				if _, err = runStack.Pop(); err != nil {
					return nil, errors.New("riskEngine: pop from runStack failed in query execute\n ")
				}
			}
		}

		if tableStr == "" {
			return nil, errors.New("riskEngine: tableStr is empty\n ")
		}

		if columnStr == "" {
			return nil, errors.New("riskEngine: columnStr is empty\n ")
		}

		executeNode := new(simpleNode)
		executeNode.Type = resultNodeType
		executeNode.Value = 0

		if len(wh) > 0 {
			mySQLStart := time.Now().UnixNano()

			jobs = append(jobs, models.Job{1, string(c.Value), columnStr, tableStr, wh})
			res, ok := QueryJob(jobs)
			//fmt.Println(res)
			if !ok {
				return nil, errors.New("riskEngine: not supported this type")
			}
			//mySQLElapsed := time.Since(mySQLStart)

			//log.Debug("DetectHandler MySQL Query Cost Time: ", (time.Now().UnixNano()-mySQLStart)/1000,"costTime")
			log.Debug("DetectHandler MySQL Query Cost Time: ", &TraceContext{
				TraceId:  TraceId,
				CostTime: (time.Now().UnixNano() - mySQLStart) / 1000,
			})
			*reason = res[0].detail

			var isThisDataInvolved = false

			for _, v := range res[0].detail {
				if v == params["prepare_id"] {
					isThisDataInvolved = true
					break
				}
			}

			if isThisDataInvolved == false {
				if BaseTime == "real_time" {
					if strings.HasPrefix(columnStr, POLYMERIZESUM) {
						// todo	sum 支持单字段和原生
						//key := columnStr[5:]
						//columnVal, ok := params[key]
						//if ok == false {
						//	return nil, errors.New("riskEngine: " + key + " is not valid\n ")
						//}
						executeNode.Value = res[0].result.(int64)
					} else if strings.HasPrefix(columnStr, POLYMERIZECOUNT) {
						executeNode.Value = res[0].result.(int64) + 1
					}
				} else {
					executeNode.Value = res[0].result
				}
			} else {
				executeNode.Value = res[0].result
			}

			//log.Info("executeNode", executeNode,"executeNode")
			//log.Info("where", wh,"where")
			log.Info("runStack", &TraceContext{
				TraceId:     TraceId,
				RunStack:    runStack,
				ExecuteNode: executeNode,
				Where:       wh,
			})

			return executeNode, nil
		} else {
			executeNode.Value = int64(0)

			return executeNode, nil
		}

	default:
		return nil, errors.New("riskEngine: not supported query type")
	}
}

// ExecuteOperatorNode execute one operator type node
func ExecuteOperatorNode(ctx context.Context, c *complexNode, runStack *Stack, params map[string]interface{}, reason *[]string) (interface{}, error) {

	opVar1, ok := (*runStack).Pop()
	if ok != nil {
		return nil, errors.New("riskEngine: Execute get op1 error\n" + ok.Error())
	}

	if opVar1.(*simpleNode).Type == fieldNodeType {

		//valueStr := reflect.ValueOf(params).Elem().FieldByName(opVar1.(*simpleNode).Value.(string))
		orderId, ok := params["subOrderId"]
		//fmt.Println(orderId)
		if ok == false {
			return nil, errors.New("riskEngine: filed orderId not exists!\n ")
		}
		rr := int(orderId.(float64))
		//fmt.Println(rr)
		orderId = strconv.Itoa(rr)
		mid := orderId.(string)
		*reason = append(*reason, mid)
		//key := changeField(opVar1.(*simpleNode).Value.(string))
		key := opVar1.(*simpleNode).Value.(string)
		valueStr, ok := params[key]
		if ok == false {
			return nil, errors.New("riskEngine: filed " + key + " not exists1!\n ")
		}
		typ := reflect.TypeOf(valueStr)
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			opVar1.(*simpleNode).Value = valueStr.(int)
		case reflect.Float32, reflect.Float64:
			opVar1.(*simpleNode).Value = valueStr.(float64)
		case reflect.String:
			opVar1.(*simpleNode).Value = valueStr.(string)
		default:
			return nil, errors.New("riskEngine: op1's Kind is Not Supported\n")
		}
	}

	opVar2, ok := (*runStack).Pop()
	if ok != nil {
		return nil, errors.New("riskEngine: Execute get op2 error\n" + ok.Error())
	}

	if opVar2.(*simpleNode).Type == fieldNodeType {
		//key := changeField(opVar2.(*simpleNode).Value.(string))
		key := opVar2.(*simpleNode).Value.(string)
		valueStr, ok := params[key]
		if ok == false {
			return nil, errors.New("riskEngine: filed " + key + " not exists!\n ")
		}
		//reflect.TypeOf(valueStr)
		typ := reflect.TypeOf(valueStr)
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			opVar2.(*simpleNode).Value = valueStr.(int)
		case reflect.Float32, reflect.Float64:
			opVar2.(*simpleNode).Value = valueStr.(float64)
		case reflect.String:
			opVar2.(*simpleNode).Value = valueStr.(string)
		default:
			return nil, errors.New("riskEngine: op2's Kind is Not Supported\n")
		}
	}

	switch c.Value {
	case numberAdd, numberMinus, numberMultiply, numberDivision:
		return ExecuteNumberOperatorOp(c, opVar1.(*simpleNode), opVar2.(*simpleNode), params)
	case compareEqual, compareNotEqual, compareGreatThan, compareGreatEqual, compareLessThan, compareLessEqual:
		return ExecuteCompareOperatorOp(c, opVar1.(*simpleNode), opVar2.(*simpleNode), params)
	case stringMinLen, stringMaxLen, stringContain, stringInWordList, stringInDict:
		return ExecuteStringOperatorOp(c, opVar1.(*simpleNode), opVar2.(*simpleNode), params)
	default:
		return nil, errors.New("riskEngine: ExecuteOperatorNode op not support:" + string(c.Value))
	}
}

// ExecuteLoginNode execute one operator type node
func ExecuteLoginNode(ctx context.Context, c *complexNode, runStack *Stack, params map[string]interface{}) (interface{}, error) {

	doLoginAnd := func(runStack *Stack) (bool, error) {
		if runStack.Len() == 0 {
			return false, nil
		}

		for _, tmpNode := range *runStack {
			if tmpNode.(*simpleNode).Value.(bool) == false {
				return false, nil
			}
			(*runStack).Pop()
		}

		return true, nil
	}

	doLoginOr := func(runStack *Stack) (bool, error) {
		if runStack.Len() == 0 {
			return false, nil
		}

		for _, tmpNode := range *runStack {
			if tmpNode.(*simpleNode).Value.(bool) == true {
				return true, nil
			}
			(*runStack).Pop()
		}

		return false, nil
	}

	doLoginNot := func(runStack *Stack) (bool, error) {
		opVar1, ok := (*runStack).Pop()
		if ok != nil {
			return false, nil
		}

		executeResult := !opVar1.(*simpleNode).Value.(bool)

		return executeResult, nil
	}

	var executeResult bool
	var err error

	switch c.Value {
	case LogicAnd:
		executeResult, err = doLoginAnd(runStack)
	case LogicOr:
		executeResult, err = doLoginOr(runStack)
	case LogicNot:
		executeResult, err = doLoginNot(runStack)
	default:
		executeResult = false
	}

	executeNode := new(simpleNode)
	executeNode.Type = boolNodeType
	executeNode.Value = executeResult

	return executeNode, err
}

// ExecuteCompareOperatorOp execute one compare operator op
func ExecuteCompareOperatorOp(op *complexNode, opVar1 *simpleNode, opVar2 *simpleNode, params map[string]interface{}) (interface{}, error) {
	doCompareEqual := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) == opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) == opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) == opVar1.Value.(string), nil
		default:
			return false, nil
		}
	}
	doCompareNotEqual := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) != opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) != opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) != opVar1.Value.(string), nil
		default:
			return false, nil
		}

	}
	doCompareGreatThan := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) > opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) > opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) > opVar1.Value.(string), nil
		default:
			return false, nil
		}
	}

	doCompareGreatEqual := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) >= opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) >= opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) >= opVar1.Value.(string), nil
		default:
			return false, nil
		}
	}
	doCompareLessThan := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) < opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) < opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) < opVar1.Value.(string), nil
		default:
			return false, nil
		}
	}
	doCompareLessEqual := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		switch opVar1.Type {
		// int value
		case integerNodeType:
			return opVar2.Value.(int64) <= opVar1.Value.(int64), nil
		// float value
		case numberNodeType:
			return opVar2.Value.(float64) <= opVar1.Value.(float64), nil
		// string value
		case stringNodeType:
			return opVar2.Value.(string) <= opVar1.Value.(string), nil
		default:
			return false, nil
		}
	}

	executeResult := false

	var err error

	switch op.Value {
	case compareEqual:
		executeResult, err = doCompareEqual(opVar1, opVar2)
	case compareNotEqual:
		executeResult, err = doCompareNotEqual(opVar1, opVar2)
	case compareGreatThan:
		executeResult, err = doCompareGreatThan(opVar1, opVar2)
	case compareGreatEqual:
		executeResult, err = doCompareGreatEqual(opVar1, opVar2)
	case compareLessThan:
		executeResult, err = doCompareLessThan(opVar1, opVar2)
	case compareLessEqual:
		executeResult, err = doCompareLessEqual(opVar1, opVar2)
	default:
		executeResult = false
	}

	executeNode := new(simpleNode)
	executeNode.Type = boolNodeType
	executeNode.Value = executeResult

	return executeNode, err
}

// ExecuteNumberOperatorOp execute one number operator op
func ExecuteNumberOperatorOp(op *complexNode, opVar1 *simpleNode, opVar2 *simpleNode, params map[string]interface{}) (*simpleNode, error) {

	doIntegerAdd := func(opVar1 *simpleNode, opVar2 *simpleNode) (int64, error) {
		return opVar2.Value.(int64) + opVar1.Value.(int64), nil
	}
	doIntegerMinus := func(opVar1 *simpleNode, opVar2 *simpleNode) (int64, error) {
		return opVar2.Value.(int64) - opVar1.Value.(int64), nil
	}
	doIntegerMultiply := func(opVar1 *simpleNode, opVar2 *simpleNode) (int64, error) {
		//强制转换到int64类型
		return opVar2.Value.(int64) * int64(opVar1.Value.(float64)*math.Pow10(0)), nil
	}
	doIntegerDiverse := func(opVar1 *simpleNode, opVar2 *simpleNode) (float64, error) {
		if opVar1.Value.(int64) == 0 {
			return 0.0, nil
		}

		return float64(opVar2.Value.(int64) / opVar1.Value.(int64)), nil
	}
	doAdd := func(opVar1 *simpleNode, opVar2 *simpleNode) (float64, error) {
		return opVar2.Value.(float64) + opVar1.Value.(float64), nil
	}
	doMinus := func(opVar1 *simpleNode, opVar2 *simpleNode) (float64, error) {
		return opVar2.Value.(float64) - opVar1.Value.(float64), nil
	}
	doMultiply := func(opVar1 *simpleNode, opVar2 *simpleNode) (float64, error) {
		return opVar2.Value.(float64) * opVar1.Value.(float64), nil
	}
	doDiverse := func(opVar1 *simpleNode, opVar2 *simpleNode) (float64, error) {
		if opVar1.Value.(float64) == 0 {
			return 0, nil
		}

		return opVar2.Value.(float64) / opVar1.Value.(float64), nil
	}

	executeNode := new(simpleNode)
	executeNode.Type = numberNodeType

	var err error

	switch op.Value {
	case numberAdd:
		if opVar1.Type == integerNodeType {
			executeNode.Value, err = doIntegerAdd(opVar1, opVar2)
			executeNode.Type = integerNodeType
		} else {
			executeNode.Value, err = doAdd(opVar1, opVar2)
		}
	case numberMinus:
		if opVar1.Type == integerNodeType {
			executeNode.Value, err = doIntegerMinus(opVar1, opVar2)
			executeNode.Type = integerNodeType
		} else {
			executeNode.Value, err = doMinus(opVar1, opVar2)
		}
	case numberMultiply:
		if opVar1.Type == integerNodeType || opVar1.Type == numberNodeType {
			executeNode.Value, err = doIntegerMultiply(opVar1, opVar2)
			executeNode.Type = integerNodeType
		} else {
			executeNode.Value, err = doMultiply(opVar1, opVar2)
		}
	case numberDivision:
		if opVar1.Type == integerNodeType {
			executeNode.Value, err = doIntegerDiverse(opVar1, opVar2)
			executeNode.Type = integerNodeType
		} else {
			executeNode.Value, err = doDiverse(opVar1, opVar2)
		}
	default:
		executeNode.Type = numberNodeType
		executeNode.Value = 0
	}

	return executeNode, err
}

// ExecuteStringOperatorOp execute one number operator op
func ExecuteStringOperatorOp(op *complexNode, opVar1 *simpleNode, opVar2 *simpleNode, params map[string]interface{}) (*simpleNode, error) {

	doMinLen := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {
		strLen := len(opVar2.String())

		return strLen < int(opVar1.Value.(float64)), nil
	}

	doMaxLen := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {
		strLen := len(opVar2.String())

		return strLen > int(opVar1.Value.(float64)), nil
	}

	doContain := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {
		return strings.Contains(opVar1.Value.(string), opVar2.Value.(string)), nil
	}

	doInWordList := func(opVar1 *simpleNode, opVar2 *simpleNode) (bool, error) {

		whereArr := strings.Split(opVar1.Value.(string), ",")

		inWordList := false

		switch opVar2.Value.(type) {
		default:
			return false, errors.New("riskEngine: op1's Kind is Not Supported in doInWordList\n")
		case int:
			opVar2.Value = strconv.Itoa(opVar2.Value.(int))
		case int64:
			opVar2.Value = strconv.FormatInt(opVar2.Value.(int64), 10)
		case int32:
			opVar2.Value = strconv.FormatInt(opVar2.Value.(int64), 10)
		case int16:
			opVar2.Value = strconv.FormatInt(opVar2.Value.(int64), 10)
		case int8:
			opVar2.Value = strconv.FormatInt(opVar2.Value.(int64), 10)
		case float32, float64:
			opVar2.Value = strconv.FormatFloat(opVar2.Value.(float64), 'f', -1, 64)
		case string:
			opVar2.Value = opVar2.Value.(string)
		}

		for _, value := range whereArr {
			if strings.EqualFold(opVar2.Value.(string), value) {
				inWordList = true
			}
		}

		return inWordList, nil
	}

	executeNode := new(simpleNode)
	executeNode.Type = boolNodeType

	var err error

	switch op.Value {
	case stringMinLen:
		executeNode.Value, err = doMinLen(opVar1, opVar2)
	case stringMaxLen:
		executeNode.Value, err = doMaxLen(opVar1, opVar2)
	case stringContain:
		executeNode.Value, err = doContain(opVar1, opVar2)
	case stringInWordList:
		executeNode.Value, err = doInWordList(opVar1, opVar2)
	default:
		executeNode.Type = boolNodeType
		executeNode.Value = false
	}

	return executeNode, err
}

// Eval the rule
func Eval(rule []byte, params map[string]interface{}) (string, bool, []string, error) {

	var haveRisk = false
	var TraceId string
	var ctx = context.TODO()

	var reason = make([]string, 0)

	if rule == nil {
		return "", haveRisk, reason, errors.New("riskEngine: rule is empty")
	}

	if params == nil {
		return "", haveRisk, reason, errors.New("riskEngine: params is nil")
	}
	cookedRule, ok := constructNodeFromString(ctx, rule)

	evalStart := time.Now().UnixNano()

	if ok != nil {
		return "", haveRisk, reason, errors.New("riskEngine: construct cooked rule failed" + ok.Error())
	}

	sign := cookedRule.Sign
	conditionRule := cookedRule.Condition
	matchRule := cookedRule.Match
	exceptionRule := cookedRule.Exception

	if conditionRule != nil {

		log.Debug("conditionRule", &TraceContext{
			TraceId:       TraceId,
			ConditionRule: conditionRule,
		})
		var conditionStack Stack

		var conditionReason = make([]string, 0)

		okCondition := conditionRule.Execute(ctx, &conditionStack, params, &conditionReason)

		if okCondition != nil {
			return sign, haveRisk, conditionReason, errors.New("riskEngine: eval condition rule failed\n" + okCondition.Error())
		}

		if conditionStack.Len() == 0 {
			return sign, haveRisk, conditionReason, errors.New("riskEngine: conditionStack is empty after Execute\n ")
		}

		var conditionRisk = true

		for _, el := range conditionStack {
			conditionRisk = conditionRisk && el.(*simpleNode).Value.(bool)
		}

		if conditionRisk == false {
			return sign, haveRisk, conditionReason, nil
		}
	}

	// check the exceptionRule
	if exceptionRule != nil {

		log.Debug("exceptionRule", exceptionRule, "exceptionRule")

		var exceptionStack Stack

		var exceptionReason = make([]string, 0)

		ok2 := exceptionRule.Execute(ctx, &exceptionStack, params, &exceptionReason)

		if ok2 != nil {
			return sign, haveRisk, exceptionReason, errors.New("riskEngine: eval exception rule failed" + ok2.Error())
		}
		var exceptionRisk = true
		if exceptionStack.Len() > 0 {
			for _, el := range exceptionStack {
				//都为true  则命中白名单
				if exceptionRisk && el.(*simpleNode).Value.(bool) {
					exceptionRisk = true
				} else {
					exceptionRisk = false
					break
				}
			}
		}
		//命中白名单，则不继续执行match规则
		if exceptionRisk == true {
			return sign, haveRisk, exceptionReason, nil
		}
	}

	//log.Debug("matchRule", matchRule,"matchRule")
	log.Debug("matchRule", &TraceContext{
		TraceId:   TraceId,
		MatchRule: matchRule,
	})
	var matchStack Stack
	ok1 := matchRule.Execute(ctx, &matchStack, params, &reason)
	if ok1 != nil {
		fmt.Println(ok1.Error())
		return sign, haveRisk, reason, errors.New("riskEngine: eval match rule failed\n" + ok1.Error())
	}

	if matchStack.Len() == 0 {
		return sign, haveRisk, reason, errors.New("riskEngine: matchStack is empty after Execute\n ")
	}

	var matchRisk = true

	for _, el := range matchStack {
		matchRisk = matchRisk && el.(*simpleNode).Value.(bool)
	}
	haveRisk = matchRisk
	//log.Debug("DetectHandler Eval Cost Time: ", (time.Now().UnixNano()-evalStart)/1000,"costTime")
	log.Debug("DetectHandler Eval Cost Time: ", &TraceContext{
		TraceId:  TraceId,
		CostTime: (time.Now().UnixNano() - evalStart) / 1000,
	})
	return sign, haveRisk, reason, nil
}
