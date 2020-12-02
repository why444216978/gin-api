package mysql

import (
	"context"
	"database/sql"
	util_err "gin-api/libraries/util/error"
	//_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	DB struct {
		masterDB *gorm.DB
		slaveDB  *gorm.DB
		Config   *Config
	}
)

func New(c *Config) (db *DB, err error) {
	db = new(DB)
	db.Config = c
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


