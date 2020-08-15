package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
	util_err "gin-frame/libraries/util/error"
	"gin-frame/libraries/xhop"

	"github.com/opentracing/opentracing-go"

	//_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	DB struct {
		IsLog    bool
		masterDB *gorm.DB
		slaveDB  *gorm.DB
		Config   *Config
	}
)

func New(c *Config) (db *DB, err error) {
	db = new(DB)
	db.Config = c
	db.IsLog = GetIsLog()
	db.masterDB, err = gorm.Open("mysql", c.Master.DSN)
	util_err.Must(err)
	db.masterDB.DB().SetMaxOpenConns(c.Master.MaxOpen)
	db.masterDB.DB().SetMaxIdleConns(c.Master.MaxIdle)
	util_err.Must(err)

	db.slaveDB, err = gorm.Open("mysql", c.Slave.DSN)
	util_err.Must(err)
	db.slaveDB.DB().SetMaxOpenConns(c.Slave.MaxOpen)
	db.slaveDB.DB().SetMaxIdleConns(c.Slave.MaxIdle)

	return
}

func (db *DB) MasterOrm() *gorm.DB {
	return db.masterDB
}

func (db *DB) SlaveOrm() *gorm.DB {
	return db.slaveDB
}

func (db *DB) MasterDB() *sql.DB {
	return db.masterDB.DB()
}

func (db *DB) SlaveDB() *sql.DB {
	return db.slaveDB.DB()
}

// MasterDBClose 释放主库的资源
func (db *DB) MasterDBClose() error {
	if db.masterDB != nil {
		err := db.masterDB.DB().Close()
		util_err.Must(err)
	}
	return nil
}

// SlaveDBClose 释放从库的资源
func (db *DB) SlaveDBClose() (err error) {
	err = db.slaveDB.DB().Close()
	util_err.Must(err)
	return nil
}

type operate int64

const (
	operateMasterExec operate = iota
	operateMasterQuery
	operateMasterQueryRow
	operateSlaveQuery
	operateSlaveQueryRow
)

var operationNames = map[operate]string{
	operateMasterExec:     "masterDBExec",
	operateMasterQuery:    "masterDBQuery",
	operateMasterQueryRow: "masterDBQueryRow",
	operateSlaveQuery:     "slaveDBQuery",
	operateSlaveQueryRow:  "slaveDBQueryRow",
}

func (db *DB) operate(ctx context.Context, op operate, query string, args ...interface{}) (i interface{}, err error) {
	var (
		parent        = opentracing.SpanFromContext(ctx)
		operationName = operationNames[op]
		span          = func() opentracing.Span {
			if parent == nil {
				return opentracing.StartSpan(operationName)
			}
			return opentracing.StartSpan(operationName, opentracing.ChildOf(parent.Context()))
		}()
		logFormat = log.LogHeaderFromContext(ctx)
		startAt   = time.Now()
		endAt     time.Time
	)

	lastModule := logFormat.Module
	lastStartTime := logFormat.StartTime
	lastEndTime := logFormat.EndTime
	lastXHop := logFormat.XHop
	defer func() {
		logFormat.Module = lastModule
		logFormat.StartTime = lastStartTime
		logFormat.EndTime = lastEndTime
		logFormat.XHop = lastXHop
	}()

	defer span.Finish()
	defer func() {
		endAt = time.Now()

		logFormat.StartTime = startAt
		logFormat.EndTime = endAt
		latencyTime := logFormat.EndTime.Sub(logFormat.StartTime).Microseconds() // 执行时间
		logFormat.LatencyTime = latencyTime
		logFormat.XHop = xhop.NewXhopNull()

		span.SetTag("error", err != nil)
		span.SetTag("db.type", "sql")
		span.SetTag("db.statement", query)
		logFormat.Module = "databus/mysql"

		if err != nil {
			db.writeError(err.Error())
			panic(err.Error())
		} else if db.IsLog == true {
			log.Infof(logFormat, "%s:[%s], params:%s, used: %d milliseconds", operationName, query,
				args, endAt.Sub(startAt).Milliseconds())
		}
	}()

	switch op {
	case operateMasterQuery:
		i, err = db.MasterDB().QueryContext(ctx, query, args...)
	case operateMasterQueryRow:
		i = db.MasterDB().QueryRowContext(ctx, query, args...)
	case operateMasterExec:
		i, err = db.MasterDB().ExecContext(ctx, query, args...)
	case operateSlaveQuery:
		i, err = db.SlaveDB().QueryContext(ctx, query, args...)
	case operateSlaveQueryRow:
		i = db.SlaveDB().QueryRowContext(ctx, query, args...)
	}
	return
}

func (db *DB) MasterDBExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	r, err := db.operate(ctx, operateMasterExec, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(sql.Result), err
}

func (db *DB) MasterDBQueryContext(ctx context.Context, query string, args ...interface{}) (result *sql.Rows, err error) {
	r, err := db.operate(ctx, operateMasterQuery, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(*sql.Rows), err
}

func (db *DB) MasterDBQueryRowContext(ctx context.Context, query string, args ...interface{}) (result *sql.Row) {
	r, _ := db.operate(ctx, operateMasterQueryRow, query, args...)
	return r.(*sql.Row)
}

func (db *DB) SlaveDBQueryContext(ctx context.Context, query string, args ...interface{}) (result *sql.Rows, err error) {
	r, err := db.operate(ctx, operateMasterQuery, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(*sql.Rows), err
}

func (db *DB) SlaveDBQueryRowContext(ctx context.Context, query string, args ...interface{}) (result *sql.Row) {
	r, _ := db.operate(ctx, operateSlaveQueryRow, query, args...)
	return r.(*sql.Row)
}

func errorsWrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func (db *DB) writeError(errMsg string) {
	errLogSection := "error"
	errorLogConfig := config.GetConfig("log", errLogSection)
	errorLogdir := errorLogConfig.Key("dir").String()

	date := time.Now().Format("2006-01-02")
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	file := errorLogdir + "/mysql/" + date + ".err"
	util.WriteWithIo(file, "["+dateTime+"]"+errMsg)
}

func GetIsLog() bool {
	cfg := config.GetConfig("log", "mysql_open")
	res, err := cfg.Key("turn").Bool()
	util.Must(err)
	return res
}
