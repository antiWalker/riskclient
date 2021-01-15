package main

import (
	"bigrisk/common"
	"bigrisk/consumer"
	"bigrisk/handlers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"
	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
	"gitlaball.nicetuan.net/wangjingnan/golib/logrus-gsr/wrapper"
	"io"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

//var logger log.Logger
func init() {
	file, _ := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{file, os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger := wrapper.NewLogger()
	logger.Logrus.SetLevel(logrus.InfoLevel)
	logger.Logrus.SetOutput(fileAndStdoutWriter)
	log.SetLogger(logger)
}

func prod() {
	consumerGroup, err := common.GetConsumerGroup()
	if err != nil {
		log.Fatal("Error creating consumer group: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumer := consumer.NewConsumer()

	doConsume := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			err := consumerGroup.Consume(ctx, common.GetTopics(), &consumer)
			if err != nil {
				log.Fatal("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}

	waitGroup := &sync.WaitGroup{}
	consumerCount, _ := common.GetConsumerCount()
	waitGroup.Add(consumerCount)

	for i := 0; i < consumerCount; i++ {
		go doConsume(waitGroup)
		log.Info("Consumer goroutine %d is up and running", i)
	}

	<-consumer.Ready
	log.Info("Consumer group is up and running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		log.Info("Consuming terminated by signal")
	case <-ctx.Done():
		log.Info("Consuming terminated by context")
	}

	cancel()
	waitGroup.Wait()

	err = consumerGroup.Close()
	if err != nil {
		log.Fatal(" Error closing consumer group: %v", err)
	}
}

func main() {
	orm.Debug = true
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
	var rules string
	//从kafka获取的order data
	params = "{\n\t\"businessAreaId\": 671,\n\t\"couponMoney\": 0,\n\t\"grouponId\": 98383,\n\t\"isNewOrder\": 0,\n\t\"mainSiteCityId\": 107,\n\t\"mainSiteCityName\": \"沈阳市\",\n\t\"mainSiteId\": 10386,\n\t\"mainSiteName\": \"沈阳市\",\n\t\"merchandiseAbbr\": \"正大 熘肉段\",\n\t\"merchandiseId\": 807950,\n\t\"merchandiseName\": \"正大 熘肉段320g\",\n\t\"merchandisePrice\": 690,\n\t\"orderId\": 450336553706083850,\n\t\"orderStatus\": 5,\n\t\"partnerId\": 271674,\n\t\"price\": 60,\n\t\"quantity\": 1,\n\t\"rebateAmount\": 69,\n\t\"siteCityId\": 107,\n\t\"siteCityName\": \"沈阳市\",\n\t\"siteId\": 10030,\n\t\"siteName\": \"沈阳市（子站）\",\n\t\"subOrderId\": 450336553706083851,\n\t\"supplyPrice\": 1523,\n\t\"ts\": 1608885401000,\n\t\"tss\": \"2020-12-25 16:36:41\",\n\t\"userId\": 118979605,\n\t\"warehouseId\": 990\n}"

	var raw = new(OrderInfo)
	if err := json.Unmarshal([]byte(params), &raw); err != nil {
		fmt.Println(err)
	}
	SiteId := "0"
	//通过子站id拼成子站场景key，然后拿着key从redis获取这个场景要过的的规则集合
	key := "RISK_FUMAOLI_SCENE_" + string(SiteId)
	rules = common.RedisGet(key)
	i, _ := raw.SubOrderId.Int64()
	ctx = context.WithValue(ctx, "TraceId", int(i))
	if rules == "" {
		//找默认的规则
		rules = common.RedisGet("RISK_FUMAOLI_SCENE_" + strconv.FormatInt(0, 10))
	}
	if rules == "" {
		fmt.Println("redis里面缓存的规则集不能为空")
		return
	}

	hit, _ := handlers.DetectHandler(params, rules, ctx)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes := hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			log.Info("命中了，命中了。", HitList, "{SubOrderId is : "+raw.SubOrderId, "{UserId     is : "+raw.UserId)
			//SubOrderId,_ := raw.SubOrderId.Int64()
			//UserId ,_    := raw.UserId.Int64()
			insertToDb(params, HitList)
		} else {
			log.Info("这个订单未命中任何策略。")
		}
	} else {
		//解析出错钉钉群报警
		fmt.Println(hit.ErrMsg)
	}
}

//如果一个订单过多条策略，则可以把这个订单下多个命中的策略批量insert。
func insertToDb(params string, HitList []handlers.StrategyResult) {
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
