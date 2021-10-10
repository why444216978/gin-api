package resource

import (
	"github.com/why444216978/gin-api/libraries/config"
	"github.com/why444216978/gin-api/libraries/etcd"
	"github.com/why444216978/gin-api/libraries/http"
	"github.com/why444216978/gin-api/libraries/logging"
	"github.com/why444216978/gin-api/libraries/orm"

	"github.com/go-redis/redis/v8"
)

var (
	Config        *config.Viper
	TestDB        *orm.Orm
	ServiceLogger *logging.Logger
	RedisCache    *redis.Client
	Etcd          *etcd.Etcd
	HTTPRPC       *http.RPC
)
