package logger

import (
	"net/http"
)

const (
	LogHeader = "Log-Id"
)

const (
	ModuleHTTP     = "HTTP"
	ModuleRPC      = "RPC"
	ModuleMySQL    = "MySQL"
	ModuleRedis    = "Redis"
	ModuleRabbitMQ = "RabbitMQ"
)

const (
	LogID      = "log_id"
	TraceID    = "trace_id"
	Header     = "header"
	Method     = "method"
	Request    = "request"
	Response   = "response"
	Code       = "code"
	ClientIP   = "client_ip"
	ClientPort = "client_port"
	ServerIP   = "server_ip"
	ServerPort = "server_port"
	Port       = "port"
	API        = "api"
	Cost       = "cost"
	Module     = "module"
	SericeName = "service_name"
	Trace      = "trace"
)

type Fields struct {
	LogID      string      `json:"log_id"`
	TraceID    string      `json:"trace_id"`
	Header     http.Header `json:"header"`
	Method     string      `json:"method"`
	Request    interface{} `json:"request"`
	Response   interface{} `json:"response"`
	Code       int         `json:"code"`
	ClientIP   string      `json:"client_ip"`
	ClientPort int         `json:"client_port"`
	ServerIP   string      `json:"server_ip"`
	ServerPort int         `json:"server_port"`
	API        string      `json:"api"`
	Cost       int64       `json:"cost"`
	Module     string      `json:"module"`
	// Trace    interface{} `json:"trace"`
}
