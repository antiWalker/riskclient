package models

import (
	"github.com/antiWalker/golib/common"
	_ "github.com/antiWalker/golib/db/orm"
	"github.com/astaxie/beego/orm"
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
	OrderTime       int64
	SubOrderId      int
	UserId          int
	MerchandiseId   int
	RuleId          int64
	Price           int
	PartnerId       int
	Quantity        int
	SupplierPrice   int
	RuleResult      string
	MerchandiseName string
	CouponMoney     int `orm:"column(CouponMoney)"`
	MerchTypeId     int `orm:"column(MerchTypeId)"`
	GrouponId       int `orm:"column(groupon_id)"`
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
	Ts               int64  `json:"ts"`
	Tss              string `json:"tss"`
	UserId           int    `json:"userId"`
	WarehouseId      int    `json:"warehouseId"`
	MerchTypeId      int    `json:"merchTypeId"`
}

// AddNegativeGrossProfitResult
//  params  kafka 收集到的数据
//  ruleId 规则id
//  对比结果 命中 还是未命中
func AddNegativeGrossProfitResult(order *Order, ruleId, ip string) (int64, error) {
	var o = orm.NewOrm()
	negativeGrossProfitResult := NegativeGrossProfitResult{}
	negativeGrossProfitResult.NegativeType = 0
	negativeGrossProfitResult.SiteId = order.MainSiteId
	negativeGrossProfitResult.OrderId = order.OrderId
	negativeGrossProfitResult.SubOrderId = order.SubOrderId
	negativeGrossProfitResult.UserId = order.UserId
	negativeGrossProfitResult.MerchandiseId = order.MerchandiseId
	negativeGrossProfitResult.MerchandiseName = order.MerchandiseName
	negativeGrossProfitResult.Price = order.Price
	negativeGrossProfitResult.SupplierPrice = order.SupplyPrice
	negativeGrossProfitResult.PartnerId = order.PartnerId
	negativeGrossProfitResult.Quantity = order.Quantity
	negativeGrossProfitResult.CouponMoney = order.CouponMoney
	negativeGrossProfitResult.MerchTypeId = order.MerchTypeId
	negativeGrossProfitResult.GrouponId = order.GrouponId
	negativeGrossProfitResult.RuleResult = ip
	if order.Ts > 0 && len(strconv.FormatInt(order.Ts, 10)) == 13 {
		negativeGrossProfitResult.OrderTime = order.Ts / 1000
	} else {
		negativeGrossProfitResult.OrderTime = order.Ts
	}

	negativeGrossProfitResult.CreateTime = time.Now().Unix()
	rule_Id, err := strconv.ParseInt(ruleId, 10, 64)
	negativeGrossProfitResult.RuleId = rule_Id
	o.Using("default")
	id, err := o.Insert(&negativeGrossProfitResult)
	if err != nil {
		common.ErrorLogger.Infof(" %v insert fail ! insert err : %v", id, err)
		return id, err
	}
	return id, err
}
