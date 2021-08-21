package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/orm"

	"github.com/go-redis/redis/v8"
)

var (
	Config        *config.Viper
	TestDB        *orm.Orm
	ServiceLogger *logging.Logger
	GormLogger    *logging.GormLogger
	RedisCache    *redis.Client
)
