package global

import (
	"encoding/json"
	"github.com/antiWalker/golib/common"
)

var RedisKey = "RISK_FUMAOLI_SCENE_"

// RuleMap 全局rule
// key : RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId)
// value : 对应的规则集合
var RuleMap = make(map[string][]string)

// RuleContentMap
//   key : RISK_FUMAOLI_RULE_1
// value : "{\"sign\":\"1\",\"match\":{\"type\":\"logic\",\"value\":\"and\",\"children\":[{\"type\":\"operator\",\"value\":\"==\",\"children\":[{\"type\":\"number\",\"value\":0},{\"type\":\"number\",\"value\":0}]},{\"type\":\"operator\",\"value\":\"<\",\"children\":[{\"type\":\"operator\",\"value\":\"-\",\"children\":[{\"type\":\"field\",\"value\":\"price\"},{\"type\":\"operator\",\"value\":\"+\",\"children\":[{\"type\":\"field\",\"value\":\"supplyPrice\"},{\"type\":\"field\",\"value\":\"couponMoney\"}]}]},{\"type\":\"number\",\"value\":0}]}]},\"exception\":{\"type\":\"logic\",\"value\":\"and\",\"children\":[{\"type\":\"operator\",\"value\":\"==\",\"children\":[{\"type\":\"number\",\"value\":0},{\"type\":\"number\",\"value\":0}]},{\"type\":\"operator\",\"value\":\"inWordList\",\"children\":[{\"type\":\"field\",\"value\":\"merchandiseId\"},{\"type\":\"string\",\"value\":\"745917,745930,745861\"}]}]}}"
var RuleContentMap = make(map[string]string)

// SetRule 设置规则
func SetRule(key, value string) {
	if value == "" {
		return
	}
	if _, ok := RuleMap[key]; ok {
		return
	}
	var ruleList []string
	if err := json.Unmarshal([]byte(value), &ruleList); err != nil {
		common.ErrorLogger.Error("rule is empty ")
	}
	RuleMap[key] = ruleList
}

// SetRuleForce  强制设置规则
func SetRuleForce(key, value string) {
	var ruleList []string
	if err := json.Unmarshal([]byte(value), &ruleList); err != nil {
		common.ErrorLogger.Error("rule is empty ")
	}
	RuleMap[key] = ruleList
}

// GetRules 获取rule集合
func GetRules(key string) []string {
	return RuleMap[key]
}

// SetRuleContent 设置规则
//   key : RISK_FUMAOLI_RULE_1
// value : "{\"sign\":\"1\",\"match\":{\"type\":\"logic\",\"value\":\"and\",\"children\":[{\"type\":\"operator\",\"value\":\"==\",\"children\":[{\"type\":\"number\",\"value\":0},{\"type\":\"number\",\"value\":0}]},{\"type\":\"operator\",\"value\":\"<\",\"children\":[{\"type\":\"operator\",\"value\":\"-\",\"children\":[{\"type\":\"field\",\"value\":\"price\"},{\"type\":\"operator\",\"value\":\"+\",\"children\":[{\"type\":\"field\",\"value\":\"supplyPrice\"},{\"type\":\"field\",\"value\":\"couponMoney\"}]}]},{\"type\":\"number\",\"value\":0}]}]},\"exception\":{\"type\":\"logic\",\"value\":\"and\",\"children\":[{\"type\":\"operator\",\"value\":\"==\",\"children\":[{\"type\":\"number\",\"value\":0},{\"type\":\"number\",\"value\":0}]},{\"type\":\"operator\",\"value\":\"inWordList\",\"children\":[{\"type\":\"field\",\"value\":\"merchandiseId\"},{\"type\":\"string\",\"value\":\"745917,745930,745861\"}]}]}}"
func SetRuleContent(key, value string) {
	if value == "" {
		return
	}
	if _, ok := RuleContentMap[key]; ok {
		return
	}
	RuleContentMap[key] = value
}

// SetRuleContentForce 强制更新
func SetRuleContentForce(key, value string) {
	if value == "" {
		return
	}
	RuleContentMap[key] = value
}

// GetRuleContent 根据key 获取value内容
func GetRuleContent(key string) string {
	return RuleContentMap[key]
}
