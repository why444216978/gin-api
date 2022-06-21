package resource

import (
	"github.com/go-redis/redis/v8"

	httpClient "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/cache"
	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/lock"
	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/queue/rabbitmq"
	etcdRegistry "github.com/why444216978/gin-api/library/registry/etcd"
)

var (
	Config        *config.Viper
	TestDB        *orm.Orm
	ServiceLogger *logger.Logger
	RedisDefault  *redis.Client
	Etcd          *etcd.Etcd
	ClientHTTP    *httpClient.RPC
	RedisLock     lock.Locker
	RedisCache    cache.Cacher
	Registrar     *etcdRegistry.EtcdRegistrar
	RabbitMQ      *rabbitmq.RabbitMQ
)
