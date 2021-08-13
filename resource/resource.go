package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/redis"

	go_redis "github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

var (
	Config       *config.Viper
	TestDB       *gorm.DB
	DefaultRedis *redis.RedisDB
	Logger       *logging.Logger
	GoRedis      *go_redis.Client
)
