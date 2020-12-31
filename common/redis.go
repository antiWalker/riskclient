package common

import (
	"github.com/go-redis/redis"
	"github.com/astaxie/beego/config"
	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
	"os"
	"time"
)
var (
	RedisCli  *redis.Client
)
func getRedisInstance() *redis.Client {
	defer func() {
		e := recover()
		if e != nil {
			log.Debug("redis recover: %s", e)
		}
	}()
	databaseConf, errDb := getConf("database")
	if errDb != nil {
		panic(errDb)
	}
	redisHost := databaseConf.String("redisHost")
	if RedisCli != nil {
		_, err := RedisCli.Ping().Result()
		if err == nil {
			return RedisCli
		}
	}
	envk8s := os.Getenv("RUNTIME_TYPE")
	var RedisAddr string
	if envk8s == "product"{
		RedisAddr = redisHost
	}else{
		RedisAddr = redisHost
	}

	RedisCli = redis.NewClient(&redis.Options{
			Addr:         RedisAddr,
			DB:           5,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     100,
			PoolTimeout:  30 * time.Second,
		})
	if _, err := RedisCli.Ping().Result(); err != nil {
		log.Debug("write redis error: %s", err)
	}
	_, err := RedisCli.Ping().Result()
	if err != nil {
		return nil
	}
	return RedisCli
}

func getConf(name string) (config.Configer, error) {
	projectRoot := os.Getenv("PROJECT_ROOT")
	confDir := ""
	if len(projectRoot) == 0 {
		confDir = "conf/" + name + ".conf"
	} else {
		confDir = projectRoot + "conf/" + name + ".conf"
	}
	databaseConf, err := config.NewConfig("ini", confDir)
	return databaseConf, err
}
func RedisGet(key string) string{
	client := getRedisInstance()

	data,err :=client.Get(key).Result()

	if err !=nil{
		data = ""
	}

	return data
}

func RedisMGet(key []string) (data []interface{},err error) {
	client := getRedisInstance()

	data,err =client.MGet(key...).Result()
	if err !=nil{
		log.Debug("MGet redis error: %s", err)
	}

	return data,err
}