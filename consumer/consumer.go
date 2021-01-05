package consumer

import (
	"bigrisk/monitor"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"riskengine/common"
	"riskengine/handlers"
	"riskengine/models"
	"strconv"
)

type Consumer struct {
	Ready chan bool
}

func NewConsumer() Consumer {
	return Consumer{
		Ready: make(chan bool),
	}
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 从kafka取params然后从redis去取rules。然后调用风控引擎模块。
	//从kafka获取的order data
	for msg := range claim.Messages() {
		fmt.Printf("Msg partition: %d, timestamp: %v, value: %s \n", msg.Partition, msg.Timestamp, string(msg.Value))
		doConsumer(string(msg.Value))
		// 手动确认消息
		session.MarkMessage(msg, "")
	}
	return nil
}

func doConsumer(params string) error {
	var raw = new(models.Order)
	if err := json.Unmarshal([]byte(params), &raw); err != nil {
		fmt.Println("json format error , ", err)
	}
	//fmt.Println("OrderInfo : ",raw)
	SiteId := raw.SiteId
	//SiteId := 10030
	//通过子站id拼成子站场景key，然后拿着key从redis获取这个场景要过的的规则集合
	key := "RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId)
	rules := common.RedisGet(key)
	fmt.Println(rules)
	if rules == "" {
		monitor.SendDingDingMessage(" 【redis里面key: RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId) + " 对应缓存的规则集不能为空，请确认数据是否异常。】")
		//fmt.Println(" RISK_FUMAOLI_SCENE_"+strconv.Itoa(SiteId)+ "redis里面缓存的规则集不能为空")
		return nil
	}
	hit, _ := handlers.DetectHandler(params, rules)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes := hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			fmt.Println("SubOrderId is : ", raw.SubOrderId)
			fmt.Println("UserId     is : ", raw.UserId)
			//SubOrderId,_ := raw.SubOrderId.Int64()
			//UserId ,_    := raw.UserId.Int64()
			fmt.Println(HitList)
			InsertToDb(params, HitList)
		}
	} else {
		//解析出错钉钉群报警
		fmt.Println(hit.ErrMsg)
	}
	return nil
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
			models.AddNegativeGrossProfitResult(params, v.Name)
		}
	}
}
