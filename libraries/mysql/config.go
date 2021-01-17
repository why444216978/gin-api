package mysql

import (
	"time"
)

type (
	Config struct {
		Master      Conn          `json:"master"`       //主库
		Slave       Conn          `json:"slave"`        //从库
		ExecTimeout time.Duration `json:"exec_timeout"` //超时打印日志
	}
	Conn struct {
		// DSN example root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true
		DSN     string `json:"dsn"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	}
)
