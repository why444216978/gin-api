package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/redis"

	"gorm.io/gorm"
)

var (
	Config       *config.Viper
	TestDB       *gorm.DB
	DefaultRedis *redis.RedisDB
	Logger       *logging.Logger
)
