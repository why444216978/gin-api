package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
)

var (
	Config       *config.Viper
	TestDB       *mysql.DB
	DefaultRedis *redis.RedisDB
	Logger       *logging.Logger
)
