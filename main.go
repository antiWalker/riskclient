package main

import (
	"bigrisk/consumer"
	"bigrisk/global"
	"bigrisk/handlers"
	"bigrisk/monitor"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlaball.nicetuan.net/wangjingnan/golib/cache/redis"
	"gitlaball.nicetuan.net/wangjingnan/golib/common"
	"gitlaball.nicetuan.net/wangjingnan/golib/mq/kafka"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func prod() {
	consumerGroup, err := kafka.GetConsumerGroup()
	if err != nil {
		common.ErrorLogger.Fatalf("Error creating consumer group: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumer := consumer.NewConsumer()

	doConsume := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			err := consumerGroup.Consume(ctx, kafka.GetConsumerTopics(), &consumer)
			if err != nil {
				common.ErrorLogger.Fatal("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}

	waitGroup := &sync.WaitGroup{}
	consumerCount, _ := kafka.GetConsumerCount()
	waitGroup.Add(consumerCount)

	for i := 0; i < consumerCount; i++ {
		go doConsume(waitGroup)
		common.InfoLogger.Infof("Consumer goroutine %d is up and running", i)
	}

	<-consumer.Ready
	common.InfoLogger.Info("Consumer group is up and running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		common.InfoLogger.Info("Consuming terminated by signal")
	case <-ctx.Done():
		common.InfoLogger.Info("Consuming terminated by context")
	}

	cancel()
	waitGroup.Wait()

	err = consumerGroup.Close()
	if err != nil {
		common.ErrorLogger.Fatalf(" Error closing consumer group: %v", err)
	}
}

func main() {
	orm.Debug = true
	go func() {
		//提供给负载均衡探活
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))

		})

		//prometheus
		http.Handle("/metrics", promhttp.Handler())
		log.Println(http.ListenAndServe("0.0.0.0:3351", nil))
	}()
	//local()
	prod()
}

func local() {
	//从订单流里面获取想要的入库信息
	type OrderInfo struct {
		OrderId    json.Number `json:"orderId"`
		SubOrderId json.Number `json:"subOrderId"`
		UserId     json.Number `json:"userId"`
		SiteId     json.Number `json:"siteId"`
	}
	ctx, _ := context.WithCancel(context.Background())
	// 从kafka取params然后从redis去取rules。然后调用风控引擎模块。
	var params string
	//从kafka获取的order data
	params = "{\n\t\"businessAreaId\": 671,\n\t\"price\": 699,\n\t\"supplyPrice\": 550,\n\t\"couponMoney\": 300,\n\t\"grouponId\": 132338,\n\t\"partnerId\": 911064,\n\t\"merchandiseId\": 816637,\n\t\"merchTypeId\": 988341,\n\t\"isNewOrder\": 0,\n\t\"mainSiteCityId\": 107,\n\t\"mainSiteCityName\": \"沈阳市\",\n\t\"mainSiteId\": 10386,\n\t\"mainSiteName\": \"沈阳市\",\n\t\"merchandiseAbbr\": \"正大 熘肉段\",\n\t\"merchandiseName\": \"正大 熘肉段320g\",\n\t\"merchandisePrice\": 690,\n\t\"orderId\": 450336553706083850,\n\t\"orderStatus\": 5,\n\t\"quantity\": 1,\n\t\"rebateAmount\": 69,\n\t\"siteCityId\": 107,\n\t\"siteCityName\": \"沈阳市\",\n\t\"siteId\": 10030,\n\t\"siteName\": \"沈阳市（子站）\",\n\t\"subOrderId\": 450336553706083851,\n\t\"ts\": 1608885401000,\n\t\"tss\": \"2020-12-25 16:36:41\",\n\t\"userId\": 118979605,\n\t\"warehouseId\": 990\n}"

	var raw = new(OrderInfo)
	if err := json.Unmarshal([]byte(params), &raw); err != nil {
		common.ErrorLogger.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(params), &data); err != nil {
		common.ErrorLogger.Fatal(err)
	}
	SiteId := "20"
	var ruleList []string
	var key string
	key = global.RedisKey + SiteId
	rules := redis.RedisGet(key)
	if rules == "" {
		//找默认的规则
		key = global.RedisKey + strconv.FormatInt(0, 10)
		ruleList = global.GetRules(key)
		if len(ruleList) == 0 {
			rules = redis.RedisGet(key)
		}
	}
	global.SetRule(key, rules)

	if len(ruleList) == 0 {
		ruleList = global.GetRules(key)
	}
	if len(ruleList) == 0 {
		monitor.SendDingDingMessage(" 【redis里面key: RISK_FUMAOLI_SCENE_" + SiteId + " 和 默认 RISK_FUMAOLI_SCENE_" + strconv.FormatInt(0, 10) + " 对应缓存的规则集不能为空，请确认数据是否异常。】")
		return
	}

	i, _ := raw.SubOrderId.Int64()
	ctx = context.WithValue(ctx, "TraceId", int(i))

	hit, _ := handlers.DetectHandler(ruleList, data, ctx)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes := hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			common.WarnLogger.Infof("TraceId : %d ,UserId : %v , 命中规则列表 :%v , ", ctx.Value("TraceId"), raw.UserId, HitList)
			insertToDb(params, HitList)
		} else {
			common.InfoLogger.Info("这个订单未命中任何策略。")
		}
	} else {
		//解析出错钉钉群报警
		common.InfoLogger.Infof("hit.ErrMsg : %v", hit.ErrMsg)
	}
}

//如果一个订单过多条策略，则可以把这个订单下多个命中的策略批量insert。
func insertToDb(params string, HitList []handlers.StrategyResult) {

	ip, _ := common.ExternalIP()
	fmt.Println(ip, "========")
	for _, v := range HitList {
		//fmt.Println(k, v)
		ruleRes := v.IsHit
		//fmt.Println(k, ruleRes)
		//把命中的策略结果insert到polardb
		if ruleRes {
			//todo imp
			//fmt.Println(ruleRes)
		}
	}
}
