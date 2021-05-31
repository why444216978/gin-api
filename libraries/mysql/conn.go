package mysql

import (
	"fmt"
	"gin-api/libraries/config"
	"strconv"
	"time"
)

type Conn struct {
	// DSN example root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true
	DSN     string `json:"dsn"`
	MaxOpen int    `json:"max_open"`
	MaxIdle int    `json:"max_idle"`
}
type Config struct {
	Master      Conn          `json:"master"`       //主库
	Slave       Conn          `json:"slave"`        //从库
	ExecTimeout time.Duration `json:"exec_timeout"` //超时打印日志
}

func GetConn(database string) (*DB, error) {
	cfg := &Config{
		Master: conn(database + "_write"),
		Slave:  conn(database + "_read"),
	}

	conn, err := newConn(cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func conn(conn string) Conn {
	return Conn{
		DSN:     getDSN(conn),
		MaxOpen: getMaxOpen(conn),
		MaxIdle: getMaxIdle(conn),
	}
}

func getExecTimeout(conn string) int64 {
	cfg := getCfg(conn)
	execTimeoutCfg := cfg["exec_timeout"].(string)
	execTimeoutCfgInt, _ := strconv.Atoi(execTimeoutCfg)
	return int64(execTimeoutCfgInt)
}

func getMaxOpen(conn string) int {
	cfg := getCfg(conn)
	maxCfg := cfg["max_open"].(string)
	maxOpen, _ := strconv.Atoi(maxCfg)
	return maxOpen
}

func getMaxIdle(conn string) int {
	cfg := getCfg(conn)
	maxCfg := cfg["max_idle"].(string)
	maxIdle, _ := strconv.Atoi(maxCfg)
	return maxIdle
}

func getDSN(conn string) string {
	cfg := getCfg(conn)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		cfg["user"].(string),
		cfg["password"].(string),
		cfg["host"].(string),
		cfg["port"].(string),
		cfg["db"].(string),
		cfg["charset"].(string))
}

func getCfg(conn string) map[string]interface{} {
	return config.GetConfigToJson("mysql", conn)
}
