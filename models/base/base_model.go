/*
 * @Descripttion:
 * @Author: weihaoyu
 */
package base

import (
	"context"
	"gin-api/libraries/config"
	"gin-api/libraries/log"
	"gin-api/libraries/mysql"
	"gin-api/libraries/util"
	util_err "gin-api/libraries/util/error"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"gopkg.in/ini.v1"
)

var cfgs map[string]*ini.Section
var dbInstance map[string]*mysql.DB

var modelInstance map[string]*BaseModel

type BaseModel struct {
	c         *gin.Context
	ctx       context.Context
	parent    opentracing.Span
	span      opentracing.Span
	logFormat *log.LogFormat

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
	util.Must(err)

	instance.Db = conn
}

func (instance *BaseModel) getExecTimeout(conn string) int64 {
	cfg := instance.getCfg(conn)
	execTimeout, err := cfg.Key("exec_timeout").Int64()
	util_err.Must(err)
	return execTimeout
}

func (instance *BaseModel) getMaxOpen(conn string) int {
	cfg := instance.getCfg(conn)
	masterNum, err := cfg.Key("max_open").Int()
	util_err.Must(err)
	return masterNum
}

func (instance *BaseModel) getMaxIdle(conn string) int {
	cfg := instance.getCfg(conn)
	masterNum, err := cfg.Key("max_idle").Int()
	util_err.Must(err)
	return masterNum
}

func (instance *BaseModel) getDSN(conn string) string {
	cfg := instance.getCfg(conn)
	dsn := cfg.Key("user").String() + ":" + cfg.Key("password").String() + "@tcp(" + cfg.Key("host").String() + ":" + cfg.Key("port").String() + ")/" + cfg.Key("db").String() + "?charset=" + cfg.Key("charset").String()
	return dsn
}

func (instance *BaseModel) getCfg(conn string) *ini.Section {
	if cfgs == nil {
		cfgs = make(map[string]*ini.Section, 30)
	}
	if cfgs[conn] == nil {
		cfgs[conn] = config.GetConfig("mysql", conn)
	}
	return cfgs[conn]
}
