package orm

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	ServiceName string
	Master      *instanceConfig
	Slave       *instanceConfig
}

type instanceConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	Charset  string
	MaxOpen  int
	MaxIdle  int
}

type Orm struct {
	*gorm.DB
	config *gorm.Config
	tracer gorm.Plugin
}

type Option func(orm *Orm)

func WithTrace(tracer gorm.Plugin) Option {
	return func(orm *Orm) {
		orm.tracer = tracer
	}
}

func WithLogger(logger logger.Interface) Option {
	return func(orm *Orm) {
		orm.config.Logger = logger
	}
}

func NewOrm(cfg *Config, opts ...Option) (orm *Orm, err error) {
	orm = &Orm{
		config: &gorm.Config{
			SkipDefaultTransaction: true,
		},
	}

	for _, o := range opts {
		o(orm)
	}

	master := mysql.Open(getDSN(cfg.Master))
	slave := mysql.Open(getDSN(cfg.Slave))

	_orm, err := gorm.Open(master, orm.config)
	if err != nil {
		err = errors.Wrap(err, "open mysql conn error：")
		return nil, err
	}

	err = _orm.Use(orm.tracer)
	if err != nil {
		return
	}

	err = _orm.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{master},
		Replicas: []gorm.Dialector{slave},
		Policy:   dbresolver.RandomPolicy{},
	}).SetMaxOpenConns(cfg.Master.MaxOpen).SetMaxIdleConns(cfg.Master.MaxIdle))
	if err != nil {
		return
	}

	orm.DB = _orm

	return
}

func (orm *Orm) UseWrite() *gorm.DB {
	return orm.Clauses(dbresolver.Write)
}

func (orm *Orm) UseRead() *gorm.DB {
	return orm.Clauses(dbresolver.Read)
}

func getDSN(cfg *instanceConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.Charset)
}
