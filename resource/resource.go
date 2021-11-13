package resource

import (
	redis_cache "github.com/why444216978/gin-api/library/cache/redis"
	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	redis_lock "github.com/why444216978/gin-api/library/lock/redis"
	"github.com/why444216978/gin-api/library/logging"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/rpc/http"

	"github.com/go-redis/redis/v8"
)

var (
	Config        *config.Viper
	TestDB        *orm.Orm
	ServiceLogger *logging.Logger
	RedisDefault  *redis.Client
	Etcd          *etcd.Etcd
	HTTPRPC       *http.RPC
	RedisLock     *redis_lock.RedisLock
	RedisCache    *redis_cache.RedisCache
)
