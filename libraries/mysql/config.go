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
		DSN     string `json:"dsn"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	}
)
