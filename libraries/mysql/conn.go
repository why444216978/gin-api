package mysql

import (
	"fmt"
	"gin-api/libraries/config"
	"strconv"
)

type BaseModel struct {
	readExecTimeout  int64
	writeExecTimeout int64
}

func GetConn(database string) (*DB, error) {
	write := database + "_write"
	read := database + "_read"
	writeDsn := getDSN(database + "_write")
	readDsn := getDSN(database + "_read")

	writeObj := Conn{
		DSN:     writeDsn,
		MaxOpen: getMaxOpen(write),
		MaxIdle: getMaxIdle(write),
	}

	readObj := Conn{
		DSN:     readDsn,
		MaxOpen: getMaxOpen(read),
		MaxIdle: getMaxOpen(read),
	}

	cfg := &Config{
		Master: writeObj,
		Slave:  readObj,
	}

	conn, err := newConn(cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
