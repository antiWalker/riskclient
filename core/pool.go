package core

import (
	"bigrisk/common"
	"bigrisk/models"
	"encoding/json"
	"errors"
	"gitlaball.nicetuan.net/wangjingnan/golib/cache/redis"
	"strconv"
	"strings"
)

//支持引擎
const (
	ENGINEMYSQL = "mysql"
	ENGINEREDIS = "redis"
)

//聚合类型
const (
	POLYMERIZECOUNT = "count"
	POLYMERIZESUM   = "sum"
	POLYMERIZEGET   = "get"
	POLYMERIZEHGET  = "hget"
)

type JobResult struct {
	job_id int64
	result interface{}
	detail []string
	error  error
}

type Engine interface {
	query(job models.Job) JobResult
}

type MysqlEngine struct {
}
type RedisEngine struct {
}

// QueryJob 引擎查询统一入口
func QueryJob(jobs []models.Job) ([]JobResult, bool) {

	var engine Engine
	var jobResults []JobResult

	for _, job := range jobs {
		//可扩展es等引擎
		if job.Engine == ENGINEMYSQL {
			engine = new(MysqlEngine)
		} else if job.Engine == ENGINEREDIS {
			engine = new(RedisEngine)
		} else {
			jobResults = append(jobResults, JobResult{job.JobId, 0, nil, errors.New("engine not exist")})
			continue
		}
		//判断是否注册了这个表
		if _, ok := models.RegStruct[job.Table]; ok {

			jobResult := engine.query(job)
			jobResults = append(jobResults, jobResult)
		} else {
			return nil, false
		}
	}

	return jobResults, true
}

// mysql引擎并发查询入口
func (mysqlEngine MysqlEngine) query(job models.Job) JobResult {
	defer func() {
		if err := recover(); err != nil {
			common.ErrorLogger.Info("捕获到了panic产生的异常 ", err)
			return
		}
	}()
	selectStruct := strings.Split(job.Select, "::")

	if len(selectStruct) != 2 {
		return JobResult{job.JobId, 0, nil, errors.New("select format error")}
	}

	var tableEngine models.TableEngine
	//路由数据表
	if job.Table != "" {
		key := job.Table
		if _, ok := models.RegStruct[key]; ok {
			engineName := models.RegStruct[key]
			tableEngine = engineName.(models.TableEngine)
		} else {
			return JobResult{job.JobId, 0, nil, errors.New("table not exist")}
		}
	} else {
		return JobResult{job.JobId, 0, nil, errors.New("table not exist")}
	}

	//路由聚合类型
	if selectStruct[0] == POLYMERIZECOUNT {
		num, detail, error := tableEngine.SpitCount(job)
		return JobResult{job.JobId, num, detail, error}
	} else if selectStruct[0] == POLYMERIZESUM {
		num, detail, error := tableEngine.SpitSum(job)
		return JobResult{job.JobId, num, detail, error}
	} else {
		return JobResult{job.JobId, 0, nil, errors.New("select type not support")}
	}
}

// redis引擎并发查询入口
func (redisEngine RedisEngine) query(job models.Job) JobResult {
	defer func() {
		if err := recover(); err != nil {
			common.ErrorLogger.Info("捕获到了panic产生的异常 ", err)
			return
		}
	}()
	selectStruct := strings.Split(job.Select, "::")
	//路由聚合类型
	if selectStruct[0] == POLYMERIZEGET {
		//merchandiseid:333:quantity:2021-01-02
		wh := job.Where
		finalKey := selectStruct[1]
		for _, e := range wh {
			if strings.Contains(finalKey, "$"+e.Column) && (e.Operator == "gt" || e.Operator == "eq") {
				finalKey = strings.ReplaceAll(finalKey, "$"+e.Column, e.Value)
			}
		}
		//get key value from redis
		common.InfoLogger.Infof("finalKey : %v ", finalKey)
		allNums, _ := strconv.ParseInt(redis.RedisGet(finalKey), 10, 64)
		return JobResult{job.JobId, allNums, []string{finalKey}, nil}
	} else if selectStruct[0] == POLYMERIZEHGET {
		// HGET sht:agbwebapp:merchandiseSale:groupon:132027:partnerId:1508368&&1082636+1309648
		wh := job.Where
		hkey := strings.Split(selectStruct[1], "&&")
		// sht:agbwebapp:merchandiseSale:groupon:132027:partnerId:1508368
		finalKey := hkey[0]
		//$merchandiseId+$merchTypeId
		field := hkey[1]
		for _, e := range wh {
			if strings.Contains(finalKey, "$"+e.Column) && (e.Operator == "gt" || e.Operator == "eq") {
				finalKey = strings.ReplaceAll(finalKey, "$"+e.Column, e.Value)
			}
			if strings.Contains(field, "$"+e.Column) && (e.Operator == "gt" || e.Operator == "eq") {
				field = strings.ReplaceAll(field, "$"+e.Column, e.Value)
			}
		}
		// json 结构 "{\"gmv\":90,\"quantity\":3,\"orderNum\":3,\"buyerNum\":3}"
		tempJson := make(map[string]int64)
		redisResult := redis.RedisHGet(finalKey, field)
		//get key value from redis
		common.InfoLogger.Infof("finalKey : %v ,filed : %v , result : %v ", finalKey, field, redisResult)
		// 防止flink 统计消费速度慢，导致获取不到数据
		if redisResult == "" {
			var n int64 = 0
			return JobResult{job.JobId, n, []string{}, errors.New("redis 取不到数据")}
		}
		if err := json.Unmarshal([]byte(redisResult), &tempJson); err != nil {
			var n int64 = 0
			common.ErrorLogger.Infof("finalKey : %v ,filed : %v , redisResult : %v, json err : %v", finalKey, field, redisResult, err)
			return JobResult{job.JobId, n, []string{}, errors.New("redis 取不到数据")}
		}
		return JobResult{job.JobId, tempJson[hkey[2]], []string{finalKey}, nil}
	} else {
		var n int64 = 0
		return JobResult{job.JobId, n, []string{}, errors.New("select type not support")}
	}
}
