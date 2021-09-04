package jaeger

const (
	fieldLogID           = "Log-Id"
	fieldTraceID         = "Trace-Id"
	fieldSpanID          = "Span-Id"
	parentSpanContextKey = "Span"
)

const (
	httpClientComponentPrefix = "HTTP-Client-"
	httpServerComponentPrefix = "HTTP-Server-"
	componentGorm             = "Gorm"
	componentRedis            = "Redis"
	componentRabbitMQ         = "RabbitMQ"
)

const (
	logFieldsRequest  = "request"
	logFieldsResponse = "response"
	logFieldsArgs     = "args"
)
