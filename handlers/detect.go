package handlers

import (
	"encoding/json"
	"riskengine/control"
	"riskengine/core"
	"sync"
	"time"

	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
)

/// 检测参数
type DetectFormV2 struct {
	/// @see json encoded `InParams` 参数
	Params string `json:"params"` // 参数
	Rules  string `json:"rules"`  // 规则
	BaseTime   string `json:"base_time"`   // 同步
}

type DetectChannel struct {
	ruleSign  string
	haveRisk  bool
	riskReason []string
	errorInfo error
}

type HitResult struct {
	IsHit bool `json:"is_hit"`
	StrategyList []StrategyResult `json:"strategy_list"`
}

type StrategyResult struct {
	Name string `json:"name"`
	IsHit  bool `json:"is_hit"`
	HitReason []string `json:"hit_reason"`
}

/// 风控检测函数
/// data => false 表示没有风险
/// data => true  表示有风险
func DetectHandler(params string,rules string) (bool,error) {
	var TraceId string

	//fmt.Println("key not found:", k)
	start :=time.Now().UnixNano()


	if core.BaseTime == "" {
		core.BaseTime = "real_time"
	}

	//log.Info("BaseTime is", core.BaseTime)
	var ruleList []interface{}
	var data interface{}

	if err := json.Unmarshal([]byte(params), &data); err != nil {
		//_ = sendResult(w, errnoParseArg, nil)
		return false,err
	}


	if err := json.Unmarshal([]byte(rules), &ruleList); err != nil {
		//_ = sendResult(w, errnoParseArg, nil)
		return false,err
	}
	//fmt.Println(ruleList)
	if len(ruleList) == 0 {
		//return sendResult(w, errnoEmptyRule, nil)
	}

	detectChannel := make(chan DetectChannel, len(ruleList))

	wg := sync.WaitGroup{}
	wg.Add(len(ruleList))
	for i := 0; i < len(ruleList); i++ {
		var ruleBytes []byte
		var err error

		if ruleBytes, err = json.Marshal(ruleList[i]); err != nil {
			//_ = sendResult(w, errnoInvalidDetectParams, nil)
			return false,err
		}

		go func(ruleBytes []byte) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err.(error).Error())
				}
			}()

			routineStart := time.Now().UnixNano()

			defer wg.Done()

			var thisDetectChannel DetectChannel

			// do detect
			if ruleSign, haveRisk, riskReason, err := control.RiskDetect(ruleBytes, data.(map[string]interface{})); err == nil {
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

			//routineElapsed := time.Since(routineStart)

			//log.Debug("DetectHandler Routine Cost Time: ",(time.Now().UnixNano()-routineStart)/1000,"costTime")
			log.Info("DetectHandler Routine Cost Time: ", &TimeContext{
				TraceId:TraceId,
				CostTime:(time.Now().UnixNano()-routineStart)/1000,
			})
		}(ruleBytes)
	}

	wg.Wait()

	close(detectChannel)

	var hitResult = new(HitResult)

	hitResult.IsHit = false

	for value := range detectChannel {
		if value.errorInfo != nil {
			//_ = sendResult(w, errnoDetectFailed, value.errorInfo.Error())
			return false,value.errorInfo
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

	//log.Info("DetectHandler Cost Time: ", (time.Now().UnixNano()-start)/1000,"costTime")
	log.Info("DetectHandler Cost Time: ", &TimeContext{
		TraceId:TraceId,
		CostTime:(time.Now().UnixNano()-start)/1000,
	})
	return hitResult.IsHit,nil
	//return sendResult(errnoSuccess, hitResult)
}
