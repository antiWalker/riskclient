package consumer

import (
	"encoding/json"
	"fmt"
	"riskengine/handlers"
)

//从订单流里面获取想要的入库信息
type OrderInfo struct {
	BusinessAreaId   json.Number `json:"businessAreaId"`
	CouponMoney      json.Number `json:"couponMoney"`
	GrouponId        json.Number `json:"grouponId"`
	IsNewOrder       json.Number `json:"isNewOrder"`
	MainSiteCityId   json.Number `json:"mainSiteCityId"`
	MainSiteCityName string      `json:"mainSiteCityName"`
	MainSiteId       json.Number `json:"mainSiteId"`
	InSiteName       string      `json:"inSiteName"`
	MerchandiseAbbr  string      `json:"merchandiseAbbr"`
	MerchandiseId    json.Number `json:"merchandiseId"`
	MerchandiseName  string      `json:"merchandiseName"`
	MerchandisePrice json.Number `json:"merchandisePrice"`
	OrderId          json.Number `json:"orderId"`
	OrderStatus      json.Number `json:"orderStatus"`
	PartnerId        json.Number `json:"partnerId"`
	Quantity         json.Number `json:"quantity"`
	RebateAmount     json.Number `json:"rebateAmount"`
	SiteCityId       json.Number `json:"siteCityId"`
	SiteCityName     string      `json:"siteCityName"`
	SiteId           int         `json:"siteId"`
	SiteName         string      `json:"siteName"`
	SubOrderId       json.Number `json:"subOrderId"`
	SupplyPrice      json.Number `json:"supplyPrice"`
	Ts               json.Number `json:"ts"`
	Tss              string      `json:"tss"`
	UserId           json.Number `json:"userId"`
	WarehouseId      json.Number `json:"warehouseId"`
}

//如果一个订单过多条策略，则可以把这个订单下多个命中的策略批量insert。
func InsertToDb(params string, HitList []handlers.StrategyResult) {
	for k, v := range HitList {
		//fmt.Println(k, v)
		ruleRes := v.IsHit
		fmt.Println(k, ruleRes)
		//把命中的策略结果insert到polardb
		if ruleRes {
			//todo imp
			//fmt.Println(ruleRes)
		}
	}
}
