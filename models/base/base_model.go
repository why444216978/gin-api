/*
 * @Descripttion:
 * @Author: weihaoyu
 */
package base

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/mysql"
	"strconv"

	"github.com/jinzhu/gorm"
	util_err "github.com/why444216978/go-util/error"
)

type BaseModel struct {
	readExecTimeout  int64
	writeExecTimeout int64
}

var cfgs map[string]interface{}
var instanceMap map[string]*mysql.DB

func init() {
	instanceMap = make(map[string]*mysql.DB, app_const.DB_NUM)
}

func (instance *BaseModel) CheckRes(dbRes *gorm.DB) error {
	if dbRes.Error != nil {
		return dbRes.Error
	}
	return nil
}

func (instance *BaseModel) GetConn(database string) *mysql.DB {
	if instanceMap[database] != nil {
		return instanceMap[database]
	}

	write := database + "_write"
	read := database + "_read"
	writeDsn := instance.getDSN(database + "_write")
	readDsn := instance.getDSN(database + "_read")

	writeObj := mysql.Conn{
		DSN:     writeDsn,
		MaxOpen: instance.getMaxOpen(write),
		MaxIdle: instance.getMaxIdle(write),
	}

	readObj := mysql.Conn{
		DSN:     readDsn,
		MaxOpen: instance.getMaxOpen(read),
		MaxIdle: instance.getMaxOpen(read),
	}

	cfg := &mysql.Config{
		Master: writeObj,
		Slave:  readObj,
	}

	instance.readExecTimeout = instance.getExecTimeout(read)
	instance.writeExecTimeout = instance.getExecTimeout(write)

	conn, err := mysql.New(cfg)
	util_err.Must(err)

	instanceMap[database] = conn

	return instanceMap[database]
}

func (instance *BaseModel) getExecTimeout(conn string) int64 {
	cfg := instance.getCfg(conn)
	execTimeoutCfg := cfg["exec_timeout"].(string)
	execTimeoutCfgInt, err := strconv.Atoi(execTimeoutCfg)
	util_err.Must(err)
	return int64(execTimeoutCfgInt)
}

func (instance *BaseModel) getMaxOpen(conn string) int {
	cfg := instance.getCfg(conn)
	maxCfg := cfg["max_open"].(string)
	maxOpen, err := strconv.Atoi(maxCfg)
	util_err.Must(err)
	return maxOpen
}

func (instance *BaseModel) getMaxIdle(conn string) int {
	cfg := instance.getCfg(conn)
	maxCfg := cfg["max_idle"].(string)
	maxIdle, err := strconv.Atoi(maxCfg)
	util_err.Must(err)
	return maxIdle
}

func (instance *BaseModel) getDSN(conn string) string {
	cfg := instance.getCfg(conn)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		cfg["user"].(string),
		cfg["password"].(string),
		cfg["host"].(string),
		cfg["port"].(string),
		cfg["db"].(string),
		cfg["charset"].(string))
}

func (instance *BaseModel) getCfg(conn string) map[string]interface{} {
	if cfgs == nil {
		cfgs = make(map[string]interface{}, 10)
	}
	if cfgs[conn] == nil {
		cfgs[conn] = config.GetConfigToJson("mysql", conn)
	}
	return cfgs[conn].(map[string]interface{})
}
