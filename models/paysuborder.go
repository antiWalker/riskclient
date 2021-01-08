package models

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "gitlaball.nicetuan.net/wangjingnan/golib/register-golang/db/orm"
	"reflect"
	"strings"
)

/*-------paysuborder---------*/

//订单表结构
type PaySubOrder struct {
	Suborderid       int64 `orm:"column(suborderid);pk"`                  // 订单id
	Merchandiseid    int64 `json:"merchandiseid"`                         // 商品id
	Orderstatus      int64 `json:"orderstatus"`                           // 状态
	Quantity         int64 `json:"quantity"`                              // 数量
	Suborderdateline int64 `json:"suborderdateline"`                      // 时间
	Price            int64 `json:"price"`                                 // 价格
	SupplyPrice      int64 `orm:"column(supplyprice)",json:"supplyprice"` // 成本价格
	Total            int64 `orm:"-",json:"total"`                         // 毛利率
}

func (u *PaySubOrder) TableName() string {
	return "paysuborder"
}

//sales-order按照条件出sum
func (paySubOrder PaySubOrder) SpitCount(job Job) (int64, []string, error) {
	orm.Debug = true
	v := func(o orm.Ormer) orm.QuerySeter {
		o.Using("slave")
		return o.QueryTable(new(PaySubOrder))
	}
	count, err := SpitCountInner(job, v)
	return count, make([]string, 0), err
}

//sales-order按照条件出sum
func (paySubOrder PaySubOrder) SpitSum(job Job) (int64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		qs := o.QueryTable(new(PaySubOrder))

		return qs
	}
	if strings.Contains(job.Select, "native") {
		return SpitSumInnerPaySubOrderNative(job, v, paySubOrder)
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
	columnName := job.Select[7:]
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

//单字段sum
func SpitSumInnerPaySubOrder(job Job, sumFunc SumFunc, tableRow interface{}) (int64, []string, error) {
	columnName := job.Select[5:]
	//转化为首字母大写
	name := getColName(columnName)
	//提取字类型
	dataType := getColumnDataType(tableRow, name)
	if dataType == reflect.Invalid {
		return 0, nil, errors.New("column " + columnName + " not allow sum")
	}
	var result int64 = 0
	var ids = make([]string, 0)
	//游标遍历条件命中的全量
	o := orm.NewOrm()
	o.Using("slave")

	var ml []PaySubOrder
	qs := sumFunc(o)

	qs = filterMerge(qs, job)
	qs.Limit(-1)
	qs.All(&ml)

	if len(ml) > 0 {
		for _, atom := range ml {
			value := getCol(atom, name)
			//累加
			result += value.(int64)
		}
	}
	fmt.Println("result:", result)
	return result, ids, nil
}

func SpitSumInnerPaySubOrderNative(job Job, sumFunc SumFunc, tableRow interface{}) (int64, []string, error) {
	fmt.Println("原生拼接sql语句")
	condition := job.Select[5:]
	columnName := strings.ReplaceAll(condition, "native:", "")
	var whereCondition bytes.Buffer
	whereCondition.WriteString("1=1 ")
	for _, where := range job.Where {
		operate, exist := OperateNative[where.Operator]
		if !exist {
			continue
		}
		whereCondition.WriteString(" and ")
		if operate == "in" {
			whereCondition.WriteString(where.Column + "  " + operate + " ( " + where.Value + " ) ")
		} else {
			whereCondition.WriteString(where.Column + " " + operate + where.Value)
		}
	}
	fmt.Println(whereCondition.String())

	o := orm.NewOrm()
	o.Using("slave")

	type Result struct {
		Total int64 `json:"total"`
	}
	var res Result
	err := o.Raw("SELECT sum(" + columnName + ") as total FROM " + job.Table + " where " + whereCondition.String() + " ").QueryRow(&res)
	if err == nil {
		fmt.Println("mysql row affected nums: ", res)
	} else {
		fmt.Println(err)
	}
	return res.Total, nil, nil
}
