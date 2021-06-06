package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"

	"github.com/opentracing/opentracing-go"
)

var (
	Config       *config.Viper
	TestDB       *mysql.DB
	DefaultRedis *redis.RedisDB
	Logger       *logging.Logger
	Tracer       opentracing.Tracer
)
