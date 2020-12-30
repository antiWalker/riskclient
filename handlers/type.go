package handlers

type (
	errnoType  int
	errMsgType string
)

const (
	errnoSuccess  errnoType = 0
	errnoFailure  errnoType = 1
	errnoParseArg errnoType = 2

	errnoInvalidOp errnoType = 100

	errnoInvalidDetectParams errnoType = 200
	errnoDetectFailed        errnoType = 201

	errnoEmptyRule            errnoType = 300
	errnoUnMatchRuleAndParams errnoType = 301
)

var errnoMsg = map[errnoType]errMsgType{
	errnoSuccess:  errMsgType("成功"),
	errnoFailure:  errMsgType("失败"),
	errnoParseArg: errMsgType("解析 HTTP 参数失败"),

	errnoInvalidOp: errMsgType("无效的 op 参数"),

	errnoInvalidDetectParams: errMsgType("无效的检测参数[aka: invalid JSON]"),
	errnoDetectFailed:        errMsgType("风控检测失败"),

	errnoEmptyRule:            errMsgType("空的规则"),
	errnoUnMatchRuleAndParams: errMsgType("规则和参数的数量不匹配"),
}

type resultType struct {
	Errno  errnoType   `json:"errno"`   // 错误码
	ErrMsg errMsgType  `json:"err_msg"` // 错误信息
	Data   interface{} `json:"data"`    // data is dependent on every handler
}

// no need for args
type SlotForm struct {
}
