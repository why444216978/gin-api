package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	redigo "gin-api/libraries/redis"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

var (
	Config       *config.Viper
	TestDB       *gorm.DB
	DefaultRedis *redigo.RedisDB
	Logger       *logging.Logger
	GoRedis      *redis.Client
)
