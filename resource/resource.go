package resource

import (
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
)

var (
	TestDB       *mysql.DB
	DefaultRedis *redis.RedisDB
)
