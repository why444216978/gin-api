package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gin-api/libraries/jaeger"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type instance struct {
	orm *gorm.DB
	db  *sql.DB
}

type Config struct {
	Master instanceConfig
	Slave  instanceConfig
}

type instanceConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DB          string
	Charset     string
	MaxOpen     int
	MaxIdle     int
	ExecTimeout int
}

type DB struct {
	master *instance
	slave  *instance
}

func NewMySQL(cfg Config) (db *DB, err error) {
	db = &DB{}

	db.master, err = getInstance(cfg.Master)
	if err != nil {
		return
	}
	db.slave, err = getInstance(cfg.Slave)
	if err != nil {
		return
	}

	return
}

func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx)
}

func (db *DB) MasterOrm() *gorm.DB {
	return db.master.orm
}

func (db *DB) SlaveOrm() *gorm.DB {
	return db.slave.orm
}

func (db *DB) MasterDB() *sql.DB {
	return db.master.db
}

func (db *DB) SlaveDB() *sql.DB {
	return db.slave.db
}

// MasterDBClose 释放主库的资源
func (db *DB) MasterDBClose() error {
	if db.master.db != nil {
		err := db.master.db.Close()
		if err != nil {
			err = errors.Wrap(err, "close master mysql conn error：")
			return err
		}
	}
	return nil
}

// SlaveDBClose 释放从库的资源
func (db *DB) SlaveDBClose() (err error) {
	err = db.slave.db.Close()
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

func getInstance(cfg instanceConfig) (conn *instance, err error) {
	db, err := sql.Open("mysql", getDSN(cfg))
	if err != nil {
		err = errors.Wrap(err, "open mysql conn error：")
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)

	orm, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "open mysql orm error：")
		return
	}
	orm.Use(jaeger.GormTrace)

	conn = &instance{
		orm: orm,
		db:  db,
	}

	return
}

func getDSN(cfg instanceConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.Charset)
}
