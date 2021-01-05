package models

import (
	"github.com/astaxie/beego/orm"
	"reflect"
	"strconv"
	"strings"
)

var Operate = map[string]string{
	"eq":  "exact",
	"gt":  "gt",
	"gte": "gte",
	"lt":  "lt",
	"lte": "lte",
	"in":  "in",
}

//支持表
const (
	TABLESALESORDER                    = "t_risk_engine_sales_order"
	TABLEHOTELORDER                    = "t_risk_engine_hotel_order"
	TABLERISKNEGATIVEGROSSPROFITRESULT = "risk_negative_gross_profit_result"
)

type TableEngine interface {
	SpitCount(job Job) (int64, []string, error)
	SpitSum(job Job) (float64, []string, error)
}

var RegStruct = make(map[string]interface{})

func init() {
	RegStruct[TABLESALESORDER] = new(SalesOrder)
	RegStruct[TABLEHOTELORDER] = new(HotelOrder)
	RegStruct[TABLERISKNEGATIVEGROSSPROFITRESULT] = new(NegativeGrossProfitResult)
	orm.RegisterModelWithPrefix("t_risk_engine_", RegStruct[TABLESALESORDER])
	orm.RegisterModelWithPrefix("t_risk_engine_", RegStruct[TABLEHOTELORDER])
	orm.RegisterModelWithPrefix("risk_", RegStruct[TABLERISKNEGATIVEGROSSPROFITRESULT])
}

type Where struct {
	Column   string
	Operator string
	Value    string
}

type Job struct {
	JobId  int64
	Engine string
	Select string
	Table  string
	Where  []Where
}

type CountFunc func(o orm.Ormer) orm.QuerySeter

type SumFunc func(o orm.Ormer) orm.QuerySeter

func getNumber(dataType reflect.Kind, value interface{}) float64 {
	if dataType == reflect.Int8 {
		return float64(value.(int8))
	} else if dataType == reflect.Int16 {
		return float64(value.(int16))
	} else if dataType == reflect.Int32 {
		return float64(value.(int32))
	} else if dataType == reflect.Int64 {
		return float64(value.(int64))
	} else if dataType == reflect.Float32 {
		return float64(value.(float32))
	} else if dataType == reflect.Float64 {
		return value.(float64)
	}
	return 0.0
}

func getString(value interface{}) string {
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case int32:
		return strconv.FormatInt(value.(int64), 10)
	case int16:
		return strconv.FormatInt(value.(int64), 10)
	case int8:
		return strconv.FormatInt(value.(int64), 10)
	case float32, float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case string:
		return value.(string)
	}

	return ""
}

// count算法
func SpitCountInner(job Job, countFunc CountFunc) (int64, error) {
	o := orm.NewOrm()

	qs := countFunc(o)
	qs = filterMerge(qs, job)

	return qs.Count()
}

func getColumnDataType(tableRow interface{}, columnName string) reflect.Kind {
	val := reflect.ValueOf(tableRow)
	s := val.FieldByName(columnName)
	return s.Kind()
}

//读取struct指定值
func getCol(tableRow interface{}, name string) interface{} {
	getType := reflect.ValueOf(tableRow).FieldByName(name)
	return getType.Interface()
}

//构造beego orm条件
func filterMerge(qs orm.QuerySeter, job Job) orm.QuerySeter {
	if len(job.Where) > 0 {
		for _, where := range job.Where {
			operate, exist := Operate[where.Operator]
			if !exist {
				continue
			}
			if operate == "in" {
				qs = qs.Filter(where.Column+"__"+operate, strings.Split(where.Value, ","))
			} else {
				qs = qs.Filter(where.Column+"__"+operate, where.Value)
			}
		}
	}
	return qs
}

//xx_xx替换为XxXx
func getColName(columnName string) string {
	var name = ""
	words := strings.Split(columnName, "_")
	if len(words) > 0 {
		for _, word := range words {
			name += wordFirstToUpper(word)
		}
	}
	return name
}

func wordFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}
