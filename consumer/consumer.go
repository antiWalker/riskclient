package consumer

import (
	"bigrisk/common"
	"bigrisk/handlers"
	"bigrisk/models"
	"bigrisk/monitor"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"gitlaball.nicetuan.net/wangjingnan/golib/cache/redis"
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
		common.InfoLogger.Infof("Msg partition: %d, timestamp: %v, value: %s \n", msg.Partition, msg.Timestamp, string(msg.Value))
		doConsumer(string(msg.Value))
		// 手动确认消息
		session.MarkMessage(msg, "")
	}
	return nil
}

func doConsumer(params string) error {
	var raw = new(models.Order)
	if err := json.Unmarshal([]byte(params), &raw); err != nil {
		common.ErrorLogger.Infof("json format error , %v ", err)
	}
	SiteId := raw.SiteId
	//SiteId := 10030
	//通过子站id拼成子站场景key，然后拿着key从redis获取这个场景要过的的规则集合
	key := "RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId)
	rules := redis.RedisGet(key)
	if rules == "" {
		//找默认的规则
		rules = redis.RedisGet("RISK_FUMAOLI_SCENE_" + strconv.FormatInt(0, 10))
	}
	if rules == "" {
		monitor.SendDingDingMessage(" 【redis里面key: RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId) + " 和 默认 RISK_FUMAOLI_SCENE_" + strconv.FormatInt(0, 10) + " 对应缓存的规则集不能为空，请确认数据是否异常。】")
		return nil
	}
	ctx, _ := context.WithCancel(context.Background())

	ctx = context.WithValue(ctx, "TraceId", raw.SubOrderId)
	hit, _ := handlers.DetectHandler(params, rules, ctx)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes := hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			common.HitLogger.Infof("TraceId : %d , Order_Info : %v , 命中规则列表 :%v ", ctx.Value("TraceId"), raw, HitList)
			InsertToDb(params, HitList)
		}
	} else {
		//解析出错钉钉群报警
		common.InfoLogger.Infof("hit.ErrMsg : %v", hit.ErrMsg)
	}
	return nil
}

//如果一个订单过多条策略，则可以把这个订单下多个命中的策略批量insert。
func InsertToDb(params string, HitList []handlers.StrategyResult) {
	for _, v := range HitList {
		//fmt.Println(k, v)
		ruleRes := v.IsHit
		//fmt.Println(k, ruleRes)
		//把命中的策略结果insert到polardb
		if ruleRes {
			//fmt.Println(ruleRes)
			id, err := models.AddNegativeGrossProfitResult(params, v.Name)
			common.SQLLogger.Infof("%v , AddNegativeGrossProfitResult : id : %v , err : %v", v.Name, id, err)
		}
	}
}
