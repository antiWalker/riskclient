package control

import (
	"bigrisk/core"
	"context"
)

/// 风控检测
func RiskDetect(rule []byte, params map[string]interface{}, context context.Context) (string, bool, []string, error) {
	return core.Eval(rule, params, context)
}
