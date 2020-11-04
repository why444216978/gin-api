package log

import (
	"context"

	"gin-api/libraries/util"
)

type logContextKey string

const ctxLogHeaderKey logContextKey = "_logHeader"

//CopyTraceInfo 拷贝context value信息，其它的全部丢弃，比如超时设置等
func CopyTraceInfo(ctx context.Context) context.Context {
	logHeader := LogHeaderFromContext(ctx)
	if logHeader == nil {
		logHeader = NewLog()
	}
	logHeader.HostIp = util.HostNamePrefix()
	return context.WithValue(context.TODO(), ctxLogHeaderKey, logHeader)
}

//LogHeaderFromContext 从context取出log.LogFormat
//logHeader 不存在会返回nil
func LogHeaderFromContext(ctx context.Context) *LogFormat {
	logHeader, _ := ctx.Value(ctxLogHeaderKey).(*LogFormat)
	return logHeader
}

//ContextWithLogHeader 挂载*LogHeader
//附加自定义的log.LogFormat
func ContextWithLogHeader(ctx context.Context, logHeader *LogFormat) context.Context {
	return context.WithValue(ctx, ctxLogHeaderKey, logHeader)
}
