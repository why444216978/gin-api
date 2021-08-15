package resource

import (
	"gin-api/libraries/config"
	"gin-api/libraries/logging"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Config     *config.Viper
	TestDB     *gorm.DB
	Logger     *logging.Logger
	RedisCache *redis.Client
)
