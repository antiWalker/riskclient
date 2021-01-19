package core

import (
	"bigrisk/common"
	"bigrisk/models"
	"errors"
	"fmt"
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
)

type JobResult struct {
	job_id int64
	result interface{}
	detail []string
	error  error
}

type Engine interface {
	query(job models.Job, result chan<- JobResult)
}

type MysqlEngine struct {
}
type RedisEngine struct {
}

/*引擎查询统一入口*/
func QueryJob(jobs []models.Job) ([]JobResult, bool) {

	var engine Engine
	var jobResults []JobResult
	ch := make(chan JobResult)

	jobNum := 0
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
		jobNum++
		//判断是否注册了这个表
		if _, ok := models.RegStruct[job.Table]; ok {
			go engine.query(job, ch)
		} else {
			return nil, false
		}
	}

	for i := 0; i < jobNum; i++ {
		jobResults = append(jobResults, <-ch)
	}
	return jobResults, true
}

// mysql引擎并发查询入口
func (mysqlEngine MysqlEngine) query(job models.Job, result chan<- JobResult) {
	defer func() {
		if err := recover(); err != nil {
			common.ErrorLogger.Info("捕获到了panic产生的异常 ", err)
			var numDefaut int64
			result <- JobResult{job.JobId, numDefaut, nil, errors.New("unknown field/column name")}
			return
		}
	}()
	selectStruct := strings.Split(job.Select, "::")

	if len(selectStruct) != 2 {
		result <- JobResult{job.JobId, 0, nil, errors.New("select format error")}
		return
	}

	var tableEngine models.TableEngine
	/*
		//路由数据表
		if job.Table == models.TABLESALESORDER {
			tableEngine = new(models.SalesOrder)
		} else if job.Table == models.TABLESTRATEGYDICTIONARY {
			tableEngine = new(models.StrategyDictionary)
		} else {
			result <- JobResult{job.JobId, 0, errors.New("table not exist")}
			return
		}
	*/
	//路由数据表
	if job.Table != "" {
		key := job.Table
		if _, ok := models.RegStruct[key]; ok {
			engineName := models.RegStruct[key]
			tableEngine = engineName.(models.TableEngine)
		} else {
			result <- JobResult{job.JobId, 0, nil, errors.New("table not exist")}
			return
		}
	} else {
		result <- JobResult{job.JobId, 0, nil, errors.New("table not exist")}
		return
	}

	//路由聚合类型
	if selectStruct[0] == POLYMERIZECOUNT {
		num, detail, error := tableEngine.SpitCount(job)
		result <- JobResult{job.JobId, num, detail, error}

	} else if selectStruct[0] == POLYMERIZESUM {
		num, detail, error := tableEngine.SpitSum(job)
		result <- JobResult{job.JobId, num, detail, error}
	} else {
		result <- JobResult{job.JobId, 0, nil, errors.New("select type not support")}
	}
}

// redis引擎并发查询入口
func (redisEngine RedisEngine) query(job models.Job, result chan<- JobResult) {
	defer func() {
		if err := recover(); err != nil {
			//if os.Getenv("RUNTIME_TYPE") == "dev" {
			//errError, _ := err.(error)
			//errString:= errError.Error()
			//log.Debug("query panic ", &TraceContext{
			//	TraceId: "",
			//	Error: errString,
			//})
			//}

			var numDefaut int64
			result <- JobResult{job.JobId, numDefaut, nil, errors.New("unknown field/column name")}
			return
		}
	}()
	selectStruct := strings.Split(job.Select, "::")
	//路由聚合类型
	if selectStruct[0] == POLYMERIZEGET {
		//merchandiseid:333:quantity:2021-01-02
		wh := job.Where
		skuKey := wh[0].Column
		skuValue := wh[0].Value
		lenTime := wh[3].Value
		dealKey := skuKey + ":" + skuValue + ":" + selectStruct[1] + ":" + lenTime
		var keyNum interface{}
		var err error
		keyNum = int64(0)
		fmt.Println(dealKey)
		//get key value from redis
		//todo imp
		var allNums []interface{}

		//var totalNumber int64
		if allNums != nil {
			keyNum = int64(1)
		} else {
			keyNum = int64(0)
		}

		result <- JobResult{job.JobId, keyNum, []string{}, err}
	} else {
		result <- JobResult{job.JobId, 0, []string{}, errors.New("select type not support")}
	}
}
