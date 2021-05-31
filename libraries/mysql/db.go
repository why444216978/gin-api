package mysql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	masterOrm *gorm.DB
	slaveOrm  *gorm.DB
	masterDB  *sql.DB
	slaveDB   *sql.DB
	Config    *Config
}

func newConn(c *Config) (db *DB, err error) {
	db = new(DB)
	db.Config = c

	db.masterDB, err = sql.Open("mysql", c.Master.DSN)
	if err != nil {
		err = errors.Wrap(err, "open master mysql conn error：")
		return nil, err
	}
	db.masterOrm, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db.masterDB,
	}), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "open master mysql orm error：")
		return
	}
	db.masterDB.SetMaxOpenConns(c.Master.MaxOpen)
	db.masterDB.SetMaxIdleConns(c.Master.MaxIdle)

	db.slaveDB, err = sql.Open("mysql", c.Slave.DSN)
	if err != nil {
		err = errors.Wrap(err, "open slave mysql error：")
		return nil, err
	}
	db.slaveOrm, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db.slaveDB,
	}), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "open slave mysql orm conn error：")
		return
	}
	db.slaveDB.SetMaxOpenConns(c.Slave.MaxOpen)
	db.slaveDB.SetMaxIdleConns(c.Slave.MaxIdle)

	return
}

func (db *DB) MasterOrm() *gorm.DB {
	return db.masterOrm
}

func (db *DB) SlaveOrm() *gorm.DB {
	return db.slaveOrm
}

func (db *DB) MasterDB() *sql.DB {
	return db.masterDB
}

func (db *DB) SlaveDB() *sql.DB {
	return db.slaveDB
}

// MasterDBClose 释放主库的资源
func (db *DB) MasterDBClose() error {
	if db.masterDB != nil {
		err := db.masterDB.Close()
		if err != nil {
			err = errors.Wrap(err, "close master mysql conn error：")
			return err
		}
	}
	return nil
}

// SlaveDBClose 释放从库的资源
func (db *DB) SlaveDBClose() (err error) {
	err = db.slaveDB.Close()
	if err != nil {
		err = errors.Wrap(err, "close slave mysql conn error：")
		return
	}
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
