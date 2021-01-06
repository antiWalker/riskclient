package handlers

import (
	"bigrisk/control/operator"
	"bigrisk/core"
	"context"
	"net/http"
)

// 一个 op 参数
type OpForm struct {
	Op string `json:"op"` // 指定的操作符号
}

// 获取指定操作符的参数 数量
func OperatorInfo(ctx context.Context, w http.ResponseWriter, args interface{}) error {
	params := args.(*OpForm)

	info := operator.OpInfo(core.AOT(params.Op))

	if info == nil {
		return sendComplexResult(w, errnoInvalidOp, "op: "+params.Op+" 是无效的参数!", nil)
	} else {
		return sendResult(errnoSuccess, info)
	}
}
