package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
	"gitlaball.nicetuan.net/wangjingnan/golib/logrus-gsr/wrapper"
	"riskengine/common"
	"riskengine/handlers"
)

//var logger log.Logger
func init() {
	logger := wrapper.NewLogger()
	logger.Logrus.SetLevel(logrus.ErrorLevel)
	log.SetLogger(logger)
}

func main() {
	//从订单流里面获取想要的入库信息
	type OrderInfo struct {
		OrderId json.Number `json:"orderId"`
		SubOrderId json.Number `json:"subOrderId"`
		UserId json.Number `json:"userId"`
		SiteId json.Number `json:"siteId"`
	}
	// 从kafka取params然后从redis去取rules。然后调用风控引擎模块。
	var params string
	var rules string
	//从kafka获取的order data
	params = "{\n\t\"businessAreaId\": 671,\n\t\"couponMoney\": 0,\n\t\"grouponId\": 98383,\n\t\"isNewOrder\": 0,\n\t\"mainSiteCityId\": 107,\n\t\"mainSiteCityName\": \"沈阳市\",\n\t\"mainSiteId\": 10386,\n\t\"mainSiteName\": \"沈阳市\",\n\t\"merchandiseAbbr\": \"正大 熘肉段\",\n\t\"merchandiseId\": 1112,\n\t\"merchandiseName\": \"正大 熘肉段320g\",\n\t\"merchandisePrice\": 690,\n\t\"orderId\": 450336553706083850,\n\t\"orderStatus\": 5,\n\t\"partnerId\": 271674,\n\t\"price\": 60,\n\t\"quantity\": 1,\n\t\"rebateAmount\": 69,\n\t\"siteCityId\": 107,\n\t\"siteCityName\": \"沈阳市\",\n\t\"siteId\": 10030,\n\t\"siteName\": \"沈阳市（子站）\",\n\t\"subOrderId\": 450336553706083851,\n\t\"supplyPrice\": 1523,\n\t\"ts\": 1608885401000,\n\t\"tss\": \"2020-12-25 16:36:41\",\n\t\"userId\": 118979605,\n\t\"warehouseId\": 990\n}"
	var raw = new(OrderInfo)
	if err := json.Unmarshal([]byte(params), &raw); err != nil {
		fmt.Println(err)
	}
	SiteId := raw.SiteId
	//通过子站id拼成子站场景key，然后拿着key从redis获取这个场景要过的的规则集合
	key := "RISK_FUMAOLI_SCENE_"+string(SiteId)
	rules = common.RedisGet(key)
	fmt.Println(rules)
	if rules ==""{
		fmt.Println("redis里面缓存的规则集不能为空")
		return
	}
	hit,_:=handlers.DetectHandler(params,rules)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes :=hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			fmt.Println("{SubOrderId is : "+raw.SubOrderId)
			fmt.Println("{UserId     is : "+raw.UserId)
			SubOrderId,_ := raw.SubOrderId.Int64()
			UserId ,_    := raw.UserId.Int64()
			insertToDb(HitList,SubOrderId,UserId)
		}
	}else{
		//解析出错钉钉群报警
		fmt.Println(hit.ErrMsg)
	}
}
func insertToDb(HitList []handlers.StrategyResult,SubOrderId int64,UserId int64)()  {
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
