package resource

import (
	"github.com/go-redis/redis/v8"

	httpClient "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/cache"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/lock"
	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/queue"
	"github.com/why444216978/gin-api/library/registry"
)

var (
	TestDB        *orm.Orm
	RedisDefault  *redis.Client
	Etcd          *etcd.Etcd
	ClientHTTP    httpClient.Client
	ServiceLogger logger.Logger
	RedisLock     lock.Locker
	RedisCache    cache.Cacher
	Registrar     registry.Registrar
	RabbitMQ      queue.Queue
)
