package handlers

import (
	"bigrisk/control"
	"bigrisk/core"
	"bigrisk/monitor"
	"context"
	"github.com/antiWalker/golib/cache/redis"
	"github.com/antiWalker/golib/common"
	"strconv"
	"sync"
	"time"
)

// DetectFormV2  检测参数
type DetectFormV2 struct {
	/// @see json encoded `InParams` 参数
	Params   string `json:"params"`    // 参数
	Rules    string `json:"rules"`     // 规则
	BaseTime string `json:"base_time"` // 同步
}

type DetectChannel struct {
	ruleSign   string
	haveRisk   bool
	riskReason []string
	errorInfo  error
}

type HitResult struct {
	IsHit        bool             `json:"is_hit"`
	StrategyList []StrategyResult `json:"strategy_list"`
}

type StrategyResult struct {
	Name      string   `json:"name"`
	IsHit     bool     `json:"is_hit"`
	HitReason []string `json:"hit_reason"`
}

// DetectHandler  风控检测函数
// data => false 表示没有风险
// data => true  表示有风险
func DetectHandler(ruleList []string, data map[string]interface{}, context context.Context) (resultType, error) {
	var TraceId = strconv.Itoa(context.Value("TraceId").(int))
	start := time.Now().UnixNano()

	if core.BaseTime == "" {
		core.BaseTime = "real_time"
	}

	//解析数组，从redis里面去
	var listKey []string
	for _, value := range ruleList {
		listKey = append(listKey, value)
	}

	keyValues, _ := redis.RedisMGet(listKey)

	detectChannel := make(chan DetectChannel, len(ruleList))

	wg := sync.WaitGroup{}
	wg.Add(len(ruleList))
	for kk, value := range keyValues {
		var ruleBytes []byte
		ruleData, ok := value.(string)
		if ok {
			ruleBytes = []byte(ruleData)
		} else {
			type rule struct {
				Ruleindex int `json:"ruleindex"`
			}
			common.ErrorLogger.Error("rule is empty ", &rule{
				Ruleindex: kk,
			})
			return makeResult(errnoEmptyRule, nil), nil
		}
		go func(ruleBytes []byte) {
			defer func() {
				if err := recover(); err != nil {
					common.ErrorLogger.Info("err : ", err)
				}
			}()

			defer wg.Done()

			var thisDetectChannel DetectChannel

			// do detect
			if ruleSign, haveRisk, riskReason, err := control.RiskDetect(ruleBytes, data, context); err == nil {
				thisDetectChannel.ruleSign = ruleSign
				thisDetectChannel.haveRisk = haveRisk
				thisDetectChannel.riskReason = riskReason
				thisDetectChannel.errorInfo = err

				detectChannel <- thisDetectChannel
			} else {
				thisDetectChannel.ruleSign = ""
				thisDetectChannel.haveRisk = false
				thisDetectChannel.riskReason = make([]string, 0)
				thisDetectChannel.errorInfo = err

				detectChannel <- thisDetectChannel
			}

		}(ruleBytes)
	}

	wg.Wait()

	close(detectChannel)

	var hitResult = new(HitResult)

	hitResult.IsHit = false

	for value := range detectChannel {
		if value.errorInfo != nil {
			common.ErrorLogger.Info(value.errorInfo)
			monitor.SendDingDingMessage(" :" + value.errorInfo.Error())
		}
		tmpStrategyResult := new(StrategyResult)
		tmpStrategyResult.Name = value.ruleSign
		tmpStrategyResult.IsHit = value.haveRisk

		if value.haveRisk == true {
			tmpStrategyResult.HitReason = value.riskReason
		} else {
			tmpStrategyResult.HitReason = make([]string, 0)
		}

		hitResult.StrategyList = append(hitResult.StrategyList, *tmpStrategyResult)

		if value.haveRisk == true {
			hitResult.IsHit = true
		}
	}

	//elapsed := time.Since(start)
	common.InfoLogger.Infof("TraceId : %v DetectHandler  Cost Time: %v", TraceId, (time.Now().UnixNano()-start)/1000)
	return makeResult(errnoSuccess, hitResult), nil
}
