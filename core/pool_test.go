package core

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"riskengine/models"
	"strconv"
	"testing"
	"time"
)

// 整形eq计算
func TestQueryJobNumberIntLogic(t *testing.T) {

	uid := time.Now().UnixNano()
	skuId := rand.Int31n(100)

	//写入条数
	num := randInt(11, 15)

	sumRes := float64(0)
	countRes := int64(0)

	a := int32(0)
	for a < num {
		a++
		err := addRow(skuId, uid)
		if err == nil {
			countRes++
			sumRes += float64(skuId)
		}
	}

	var jobs []models.Job
	var wh []models.Where

	wh = append(wh, models.Where{"uid", "eq", strconv.FormatInt(uid, 10)})
	wh = append(wh, models.Where{"sku_id", "eq", fmt.Sprintf("%d", skuId)})

	jobs = append(jobs, models.Job{1, "mysql", "count::id", "sales_order", wh})
	jobs = append(jobs, models.Job{2, "mysql", "sum::sku_id", "sales_order", wh})

	res ,_:= QueryJob(jobs)

	for _, atom := range res {
		if atom.job_id == 1 {
			if atom.result.(int64) != countRes {
				t.Errorf("count error,the answer is %d,actual is %d", atom.result, countRes)
			}
		}

		if atom.job_id == 2 {

			if atom.result.(float64) != float64(sumRes) {
				t.Errorf("count error,the answer is %d,actual is %d", atom.result, countRes)
			}
		}
	}
}

func randInt(min int32, max int32) int32 {
	return min + rand.Int31n(max-min)
}

func addRow(skuId int32, uid int64) error {
	o := orm.NewOrm()
	job := new(models.SalesOrder)

	job.SkuId = skuId
	job.Uid = uid

	_, insert_err := o.Insert(job)

	return insert_err
}
