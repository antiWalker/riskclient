package common

import (
	"encoding/json"
)

var RedisKey = "RISK_FUMAOLI_SCENE_"

// RuleMap 全局rule
// key : RISK_FUMAOLI_SCENE_" + strconv.Itoa(SiteId)
// value : 对应的规则集合
var RuleMap = make(map[string][]string)

// SetRule 设置规则
func SetRule(key, value string) {
	var ruleList []string
	if err := json.Unmarshal([]byte(value), &ruleList); err != nil {
		ErrorLogger.Error("rule is empty ")
	}
	RuleMap[key] = ruleList
}

// GetRules 获取rule集合
func GetRules(key string) []string {
	return RuleMap[key]
}
