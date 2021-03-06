package consumer

import (
	"bigrisk/global"
	"bigrisk/handlers"
	"bigrisk/models"
	"bigrisk/monitor"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/antiWalker/golib/cache/redis"
	"github.com/antiWalker/golib/common"
	"os"
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

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(params), &data); err != nil {
		common.ErrorLogger.Infof("err : %v", err)
	}

	SiteId := raw.SiteId
	//SiteId := 10030
	//通过子站id拼成子站场景key，然后拿着key从redis获取这个场景要过的的规则集合
	var ruleList []string
	var key string
	key = global.RedisKey + strconv.Itoa(SiteId)
	rules := redis.RedisGet(key)
	if rules == "" {
		//找默认的规则
		key = global.RedisKey + strconv.FormatInt(0, 10)
		ruleList = global.GetRules(key)
		if len(ruleList) == 0 {
			common.WarnLogger.Info("从redis中取规则集")
			rules = redis.RedisGet(key)
		}
	}
	global.SetRule(key, rules)

	if len(ruleList) == 0 {
		ruleList = global.GetRules(key)
	}
	if len(ruleList) == 0 {
		ip, _ := common.ExternalIP()
		hostname, _ := os.Hostname()
		monitor.SendDingDingMessage("ip : " + ip.String() + " [ " + hostname + " ]【redis里面key: RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId) + " 和 默认 RISK_FUMAOLI_SCENE_" + strconv.FormatInt(0, 10) + " 对应缓存的规则集不能为空，请确认数据是否异常。】")
		return nil
	}

	ctx, _ := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "TraceId", raw.SubOrderId)
	hit, _ := handlers.DetectHandler(ruleList, data, ctx)
	//解析结果
	Errno := hit.Errno
	if Errno == 0 {
		HRes := hit.Data.(*handlers.HitResult)
		IsHit := HRes.IsHit
		HitList := HRes.StrategyList
		//命中后再做log到db的操作。
		if IsHit == true {
			common.WarnLogger.Infof("TraceId : %d , Order_Info : %v , 命中规则列表 :%v ", ctx.Value("TraceId"), raw, HitList)
			InsertToDb(ctx, raw, HitList)
		}
	} else {
		//解析出错钉钉群报警
		common.InfoLogger.Infof("hit.ErrMsg : %v", hit.ErrMsg)
	}
	return nil
}

// InsertToDb 如果一个订单过多条策略，则可以把这个订单下多个命中的策略批量insert。
func InsertToDb(ctx context.Context, order *models.Order, HitList []handlers.StrategyResult) {
	ip, _ := common.ExternalIP()
	for _, v := range HitList {
		ruleRes := v.IsHit
		if ruleRes {
			//fmt.Println(ruleRes)
			id, err := models.AddNegativeGrossProfitResult(order, v.Name, ip.String())
			common.WarnLogger.Infof("TraceId : %d ,StrategyResult : %v , AddNegativeGrossProfitResult : id : %v , err : %v", ctx.Value("TraceId"), v, id, err)
		}
	}
}
