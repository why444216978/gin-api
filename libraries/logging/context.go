package logging

import (
	"context"
)

type logContextKey string

const ctxLogHeaderKey logContextKey = "_logHeader"

//LogHeaderFromContext 从context取出log.LogFormat
//logHeader 不存在会返回nil
func LogHeaderFromContext(ctx context.Context) *LogHeader {
	logHeader, _ := ctx.Value(ctxLogHeaderKey).(*LogHeader)
	return logHeader
}

//ContextWithLogHeader 挂载*LogHeader
//附加自定义的log.LogFormat
func ContextWithLogHeader(ctx context.Context, logHeader *LogHeader) context.Context {
	return context.WithValue(ctx, ctxLogHeaderKey, logHeader)
}
