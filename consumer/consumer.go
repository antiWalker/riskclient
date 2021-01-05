package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"riskengine/common"
	"riskengine/handlers"
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
		fmt.Printf("Msg partition: %d, timestamp: %v, value: %s", msg.Partition, msg.Timestamp, string(msg.Value))
		doConsumer(string(msg.Value))
		// 手动确认消息
		session.MarkMessage(msg, "")
	}
	return nil
}

func doConsumer(params string) error {
	var raw = new(OrderInfo)
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
		fmt.Println("RISK_FUMAOLI_SCENE_"+strconv.Itoa(SiteId), "redis里面缓存的规则集不能为空")
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
			fmt.Println("{SubOrderId is : " + raw.SubOrderId)
			fmt.Println("{UserId     is : " + raw.UserId)
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
