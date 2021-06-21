package jaeger

const (
	fieldLogID           = "Log-Id"
	fieldTraceID         = "Trace-Id"
	fieldSpanID          = "Span-Id"
	parentSpanContextKey = "Span"
)

const (
	operationTypeHTTP     = "HTTP"
	operationTypeGorm     = "Gorm"
	operationTypeRedis    = "Redis"
	operationTypeRabbitMQ = "RabbitMQ"
)

const (
	logFieldsRequest  = "request"
	logFieldsResponse = "response"
	logFieldsArgs     = "args"
)
