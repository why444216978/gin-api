/*
 * @Descripttion:
 * @Author: weihaoyu
 */
package base

import (
	"context"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	util_err "github.com/why444216978/go-util/error"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"strconv"
)

var cfgs map[string]interface{}
var dbInstance map[string]*mysql.DB

var modelInstance map[string]*BaseModel

type BaseModel struct {
	c         *gin.Context
	ctx       context.Context
	parent    opentracing.Span
	span      opentracing.Span
	logFormat *logging.LogHeader

	readExecTimeout  int64
	writeExecTimeout int64

	Db *mysql.DB
}

func (instance *BaseModel) CheckRes(dbRes *gorm.DB) {
	if dbRes.Error != nil {
		panic(dbRes.Error)
	}
}

func (instance *BaseModel) GetConn(database string) {
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

	instance.Db = conn
}

func (instance *BaseModel) getExecTimeout(conn string) int64 {
	cfg := instance.getCfg(conn)
	execTimeoutCfg := cfg["exec_timeout"].(string)
	execTimeoutCfgInt,err := strconv.Atoi(execTimeoutCfg)
	util_err.Must(err)
	return int64(execTimeoutCfgInt)
}

func (instance *BaseModel) getMaxOpen(conn string) int {
	cfg := instance.getCfg(conn)
	maxCfg := cfg["max_open"].(string)
	maxOpen,err := strconv.Atoi(maxCfg)
	util_err.Must(err)
	return maxOpen
}

func (instance *BaseModel) getMaxIdle(conn string) int {
	cfg := instance.getCfg(conn)
	maxCfg := cfg["max_idle"].(string)
	maxIdle,err := strconv.Atoi(maxCfg)
	util_err.Must(err)
	return maxIdle
}

func (instance *BaseModel) getDSN(conn string) string {
	cfg := instance.getCfg(conn)
	dsn := cfg["user"].(string) + ":" + cfg["password"].(string) + "@tcp(" + cfg["host"].(string) + ":" + cfg["port"].(string) + ")/" + cfg["db"].(string) + "?charset=" + cfg["charset"].(string)
	return dsn
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
