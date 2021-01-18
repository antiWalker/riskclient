package common

import (
	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var InfoLogger *zap.SugaredLogger
var DebugLogger *zap.SugaredLogger
var ErrorLogger *zap.SugaredLogger
var WarnLogger *zap.SugaredLogger
var HitLogger *zap.SugaredLogger
var SQLLogger *zap.SugaredLogger

func init() {
	InfoLogger = log.NewLogger("./logs/info.log", zapcore.InfoLevel, 500, 30, 7, true, "riskclient").Sugar()
	DebugLogger = log.NewLogger("./logs/debug.log", zapcore.DebugLevel, 500, 30, 7, true, "riskclient").Sugar()
	ErrorLogger = log.NewLogger("./logs/error.log", zapcore.ErrorLevel, 500, 30, 7, true, "riskclient").Sugar()
	WarnLogger = log.NewLogger("./logs/warn.log", zapcore.WarnLevel, 500, 30, 7, true, "riskclient").Sugar()
	HitLogger = log.NewLogger("./logs/hit.log", zapcore.InfoLevel, 500, 30, 7, true, "riskclient").Sugar()
	SQLLogger = log.NewLogger("./logs/sql.log", zapcore.InfoLevel, 500, 30, 7, true, "riskclient").Sugar()
}
