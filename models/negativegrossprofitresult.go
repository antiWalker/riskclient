package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "gitlaball.nicetuan.net/wangjingnan/golib/register-golang/db/orm"
	"strconv"
	"time"
)

type NegativeGrossProfitResult struct {
	NegativeId      int `orm:"pk;auto;column(negative_id)"` //表示设置为主键并且自增，列名为id
	NegativeType    int
	SiteId          int
	CreateTime      int64
	UpdateTime      int64
	OrderId         int
	SubOrderId      int
	UserId          int
	MerchandiseId   int
	RuleId          int64
	Price           int
	PartnerId       int
	SupplierPrice   int
	RuleResult      string
	MerchandiseName string
}

type Order struct {
	BusinessAreaId   int    `json:"businessAreaId"`
	CouponMoney      int    `json:"couponMoney"`
	GrouponId        int    `json:"grouponId"`
	IsNewOrder       int    `json:"isNewOrder"`
	MainSiteCityId   int    `json:"mainSiteCityId"`
	MainSiteCityName string `json:"mainSiteCityName"`
	MainSiteId       int    `json:"mainSiteId"`
	MainSiteName     string `json:"mainSiteName"`
	MerchandiseAbbr  string `json:"merchandiseAbbr"`
	MerchandiseId    int    `json:"merchandiseId"`
	MerchandiseName  string `json:"merchandiseName"`
	MerchandisePrice int    `json:"merchandisePrice"`
	OrderId          int    `json:"orderId"`
	OrderStatus      int    `json:"orderStatus"`
	PartnerId        int    `json:"partnerId"`
	Price            int    `json:"price"`
	Quantity         int    `json:"quantity"`
	RebateAmount     int    `json:"rebateAmount"`
	SiteCityId       int    `json:"siteCityId"`
	SiteCityName     string `json:"siteCityName"`
	SiteName         string `json:"siteName"`
	SiteId           int    `json:"siteId"`
	SubOrderId       int    `json:"subOrderId"`
	SupplyPrice      int    `json:"supplyPrice"`
	Ts               int    `json:"ts"`
	Tcs              string `json:"tcs"`
	UserId           int    `json:"userId"`
	WarehouseId      int    `json:"warehouseId"`
}

/**
  params  kafka 收集到的数据
  ruleId 规则id
  对比结果 命中 还是未命中
*/
func AddNegativeGrossProfitResult(params string, ruleId string) (int64, error) {
	var order Order
	if err := json.Unmarshal([]byte(params), &order); err == nil {
		fmt.Println(order)
		var o = orm.NewOrm()
		negativeGrossProfitResult := NegativeGrossProfitResult{}
		negativeGrossProfitResult.NegativeType = 0
		negativeGrossProfitResult.SiteId = order.SiteId
		negativeGrossProfitResult.OrderId = order.OrderId
		negativeGrossProfitResult.SubOrderId = order.SubOrderId
		negativeGrossProfitResult.UserId = order.UserId
		negativeGrossProfitResult.MerchandiseId = order.MerchandiseId
		negativeGrossProfitResult.MerchandiseName = order.MerchandiseName
		negativeGrossProfitResult.Price = order.Price
		negativeGrossProfitResult.SupplierPrice = order.SupplyPrice
		negativeGrossProfitResult.PartnerId = order.PartnerId
		negativeGrossProfitResult.CreateTime = time.Now().Unix()
		rule_Id, err := strconv.ParseInt(ruleId, 10, 64)
		negativeGrossProfitResult.RuleId = rule_Id
		id, err := o.Insert(&negativeGrossProfitResult)
		if err != nil {
			fmt.Println("insert err :", err)
			return id, err
		}

		return 0, err
	} else {
		fmt.Println(err)
		return 0, err
	}
}
