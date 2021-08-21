package logging

import (
	"net/http"
)

const (
	LogHeader = "Log-Id"
)

type MODULE string

const (
	ModuleHTTP     MODULE = "HTTP"
	ModuleRPC      MODULE = "RPC"
	ModuleMySQL    MODULE = "MySQL"
	ModuleRedis    MODULE = "Redis"
	ModuleRabbitMQ MODULE = "RabbitMQ"
)

const (
	LogID    = "log_id"
	TraceID  = "trace_id"
	Header   = "header"
	Method   = "method"
	Request  = "request"
	Response = "response"
	Code     = "code"
	CallerIP = "caller_ip"
	HostIP   = "host_ip"
	Port     = "port"
	API      = "api"
	Cost     = "cost"
	Module   = "module"
	Trace    = "trace"
)

type Fields struct {
	LogID    string      `json:"log_id"`
	TraceID  string      `json:"trace_id"`
	Header   http.Header `json:"header"`
	Method   string      `json:"method"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	Code     int         `json:"code"`
	CallerIP string      `json:"caller_ip"`
	HostIP   string      `json:"host_ip"`
	API      string      `json:"api"`
	Cost     int64       `json:"cost"`
	Module   MODULE      `json:"module"`
	// Trace    interface{} `json:"trace"`
}
