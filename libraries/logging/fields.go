package logging

import (
	"net/http"
)

const (
	LOG_FIELD = "Log-Id"
)

type MODULE string

const (
	MODULE_HTTP     = "HTTP"
	MODULE_RPC      = "RPC"
	MODULE_MYSQL    = "MySQL"
	MODULE_REDIS    = "Redis"
	MODULE_RabbitMQ = "RabbitMQ"
)

type Common struct {
	LogID string `json:"log_id"`
}

type Fields struct {
	LogID    string      `json:"log_id"`
	TraceID  string      `json:"trace_id"`
	SpanID   string      `json:"span_id"`
	Header   http.Header `json:"header"`
	Method   string      `json:"method"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	Code     int         `json:"code"`
	CallerIP string      `json:"caller_ip"`
	HostIP   string      `json:"host_ip"`
	Port     int         `json:"port"`
	API      string      `json:"api"`
	Cost     int64       `json:"cost"`
	Module   string      `json:"module"`
	Trace    interface{} `json:"trace"`
}
