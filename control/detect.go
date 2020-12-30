package control

import (
	"riskengine/core"
)

/// 风控检测
func RiskDetect(rule []byte, params map[string]interface{}) (string, bool, []string, error) {
	return core.Eval(rule, params)
}
