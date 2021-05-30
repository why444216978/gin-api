package resource

import (
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
)

var (
	TestDB       *mysql.DB
	DefaultRedis *redis.RedisDB
	Logger       *logging.Logger
)
