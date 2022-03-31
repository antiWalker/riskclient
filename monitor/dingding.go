package monitor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antiWalker/golib/cache/redis"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/nacos"
	"net/http"
	"time"
)

func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Sign 获取加密后的url地址
func Sign() string {
	secret := nacos.GetConfigMap()["dingding.security"]
	webhook := nacos.GetConfigMap()["dingding.webHook"]
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	return url
}

// SendDingDingMessage 发送消息
func SendDingDingMessage(contentData string) bool {
	if checkLimit() {
		common.WarnLogger.Info("钉钉消息超过限制，不发送。")
		return false
	}
	var atMap = make(map[string]string)
	atMap["isAtAll"] = "true"
	atMap["msgtype"] = "text"
	content, data := make(map[string]string), make(map[string]interface{})
	content["content"] = nacos.GetConfigMap()["dingding.keyWords"] + contentData
	data["msgtype"] = "text"
	data["text"] = content
	data["at"] = atMap
	b, _ := json.Marshal(data)

	resp, err := http.Post(Sign(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		common.ErrorLogger.Info(err)
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return true
}

//check 检查消息发送频率是否超过限制：每个机器人每分钟最多发送20条。如果超过20条，会限流10分钟。
func checkLimit() bool {
	//发送前设置一个缓存一分钟的自增数
	value := redis.RedisIncrEx(nacos.GetConfigMap()["dingding.riskclientId"], 60*time.Second)
	if value > 20 {
		return true
	}
	return false
}
