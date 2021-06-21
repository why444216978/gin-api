package mysql

import (
	"context"
	"fmt"
	"gin-api/libraries/jaeger"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

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

func NewMySQL(cfg Config) (orm *gorm.DB, err error) {
	master := mysql.Open(getDSN(cfg.Master))
	slave := mysql.Open(getDSN(cfg.Slave))

	orm, err = gorm.Open(master, &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "open mysql conn error：")
		return nil, err
	}

	orm.Use(jaeger.GormTrace)

	orm.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{master},
		Replicas: []gorm.Dialector{slave},
		Policy:   dbresolver.RandomPolicy{},
	}).SetMaxOpenConns(cfg.Master.MaxOpen).SetMaxIdleConns(cfg.Master.MaxIdle))

	return
}

func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx)
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
