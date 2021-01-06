package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "gitlaball.nicetuan.net/wangjingnan/golib/register-golang/db/orm"
	"reflect"
)

/*-------paysuborder---------*/

//订单表结构
type PaySubOrder struct {
	Suborderid             int64   `orm:"column(suborderid);pk"`      // 订单id
	Merchandiseid          int64   `json:"merchandiseid"`             // 商品id
	Orderstatus            int64   `json:"orderstatus"`               // 状态
	Quantity               int64   `json:"quantity"`                  // 数量
	Suborderdateline       int64   `json:"suborderdateline"`          // 时间
}

//sales-order按照条件出sum
func (paySubOrder PaySubOrder) SpitCount(job Job) (int64, []string, error) {
	orm.Debug=true
	v := func(o orm.Ormer) orm.QuerySeter {
		return o.QueryTable(new(PaySubOrder))
	}
	count, err := SpitCountInner(job, v)
	return count, make([]string, 0), err
}

//sales-order按照条件出sum
func (paySubOrder PaySubOrder) SpitSum(job Job) (float64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		qs := o.QueryTable(new(PaySubOrder))

		return qs
	}
	return SpitSumInnerPaySubOrder(job, v, paySubOrder)
}

func SplitCountInnerPaySubOrder(job Job, countFunc CountFunc, tableRow interface{}) (int64, []string, error) {
	var id int64 = 0
	/*
		selectStruct := strings.Split(job.Select, "::")
		if len(selectStruct) != 2 {
			return 0, nil, errors.New("select format error")
		}

		var columnName = selectStruct[1]
	*/
	columnName :=job.Select[7 : ]
	//转化为首字母大写
	name := getColName(columnName)
	//提取字类型
	dataType := getColumnDataType(tableRow, name)
	if dataType == reflect.Invalid {
		fmt.Println("column " + columnName + " not allow count")
		return 0, nil, errors.New("column " + columnName + " not allow count")
	}

	result := make(map[string]interface{})

	var ids = make([]string, 0)

	//游标遍历条件命中的全量
	for true {
		o := orm.NewOrm()

		var ml []PaySubOrder
		qs := countFunc(o)
		qs = qs.Filter("suborderid__gt", id)
		qs = filterMerge(qs, job)
		qs.Limit(100).OrderBy("suborderid")
		qs.All(&ml)
		if len(ml) > 0 {
			for _, atom := range ml {
				id = atom.Suborderid

				value := getCol(atom, name)

				key := getString(value)

				if _, exists := result[key]; exists == false {
					result[key] = value

					ids = append(ids, string(atom.Suborderid))
				}
			}
		} else {
			break
		}
	}

	return int64(len(result)), ids, nil
}

func SpitSumInnerPaySubOrder(job Job, sumFunc SumFunc, tableRow interface{}) (float64, []string, error) {
	var id int64 = 0
	/*
		selectStruct := strings.Split(job.Select, "::")
		if len(selectStruct) != 2 {
			return 0, nil, errors.New("select format error")
		}

		var columnName = selectStruct[1]
	*/
	columnName :=job.Select[5 : ]
	//转化为首字母大写
	name := getColName(columnName)

	//提取字类型
	dataType := getColumnDataType(tableRow, name)
	if dataType == reflect.Invalid {
		return 0, nil, errors.New("column " + columnName + " not allow sum")
	}

	var result = 0.0
	var ids = make([]string, 0)

	//游标遍历条件命中的全量
	for true {
		o := orm.NewOrm()

		var ml []SalesOrder
		qs := sumFunc(o)

		qs = qs.Filter("id__gt", id)

		qs = filterMerge(qs, job)
		qs.Limit(100).OrderBy("id")
		qs.All(&ml)

		if len(ml) > 0 {
			for _, atom := range ml {
				id = atom.Id

				value := getCol(atom, name)

				//累加
				result += getNumber(dataType, value)

				ids = append(ids, atom.PrepareId)

			}
		} else {
			break
		}
	}
	return result, ids, nil
}
