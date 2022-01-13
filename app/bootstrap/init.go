package bootstrap

import (
	"context"
	"flag"
	"log"
	"path"
	"path/filepath"
	"strings"

	appConfig "github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/client/codec"
	httpClient "github.com/why444216978/gin-api/client/http"
	redisCache "github.com/why444216978/gin-api/library/cache/redis"
	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/jaeger"
	jaegerGorm "github.com/why444216978/gin-api/library/jaeger/gorm"
	jaegerRedis "github.com/why444216978/gin-api/library/jaeger/redis"
	redisLock "github.com/why444216978/gin-api/library/lock/redis"
	"github.com/why444216978/gin-api/library/logger"
	loggerGorm "github.com/why444216978/gin-api/library/logger/gorm"
	loggerRedis "github.com/why444216978/gin-api/library/logger/redis"
	loggerRPC "github.com/why444216978/gin-api/library/logger/rpc"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/redis"
	"github.com/why444216978/gin-api/library/registry"
	registryEtcd "github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/library/servicer"
)

var (
	envFlag = flag.String("env", "dev", "config path")
)

var envMap = map[string]struct{}{
	"dev":      struct{}{},
	"liantiao": struct{}{},
	"qa":       struct{}{},
	"online":   struct{}{},
}

var (
	env      string
	confPath string
)

func initResource(ctx context.Context) {
	initConfig()
	initApp()
	initLogger()
	initMysql("test_mysql")
	initRedis("default_redis")
	initJaeger()
	initEtcd()
	initServices(ctx)
	initClientHTTP()
	initLock()
	initCache()
}

func initConfig() {
	env = *envFlag
	log.Println("The environment is :" + env)

	if _, ok := envMap[env]; !ok {
		panic(env + " error")
	}

	confPath = "conf/" + env

	var err error
	resource.Config = config.InitConfig(confPath, "toml")
	if err != nil {
		panic(err)
	}
}

func initApp() {
	if err := resource.Config.ReadConfig("app", "toml", &appConfig.App); err != nil {
		panic(err)
	}
}

func initLogger() {
	var err error
	cfg := &logger.Config{}

	if err = resource.Config.ReadConfig("log/service", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.ServiceLogger, err = logger.NewLogger(cfg,
		logger.WithModule(logger.ModuleHTTP),
		logger.WithServiceName(appConfig.App.AppName),
	)
	if err != nil {
		panic(err)
	}

	RegisterCloseFunc(resource.ServiceLogger.Sync())
}

func initMysql(db string) {
	var err error
	cfg := &orm.Config{}
	logCfg := &loggerGorm.GormConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}

	if err = resource.Config.ReadConfig("log/gorm", "toml", logCfg); err != nil {
		panic(err)
	}

	logCfg.ServiceName = cfg.ServiceName
	gormLogger, err := loggerGorm.NewGorm(logCfg)
	if err != nil {
		panic(err)
	}

	resource.TestDB, err = orm.NewOrm(cfg,
		orm.WithTrace(jaegerGorm.GormTrace),
		orm.WithLogger(gormLogger),
	)
	if err != nil {
		panic(err)
	}
}

func initRedis(db string) {
	var err error
	cfg := &redis.Config{}
	logCfg := &loggerRedis.RedisConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}
	if err = resource.Config.ReadConfig("log/redis", "toml", &logCfg); err != nil {
		panic(err)
	}

	logCfg.ServiceName = cfg.ServiceName
	logCfg.Host = cfg.Host
	logCfg.Port = cfg.Port

	logger, err := loggerRedis.NewRedisLogger(logCfg)
	if err != nil {
		panic(err)
	}

	rc := redis.NewClient(cfg)
	rc.AddHook(jaegerRedis.NewJaegerHook())
	rc.AddHook(logger)
	resource.RedisDefault = rc
}

func initLock() {
	var err error
	resource.RedisLock, err = redisLock.New(resource.RedisDefault)
	if err != nil {
		panic(err)
	}
}

func initCache() {
	var err error

	resource.RedisCache, err = redisCache.New(resource.RedisDefault, resource.RedisLock)
	if err != nil {
		panic(err)
	}
}

func initJaeger() {
	var err error
	cfg := &jaeger.Config{}

	if err = resource.Config.ReadConfig("jaeger", "toml", cfg); err != nil {
		panic(err)
	}

	_, _, err = jaeger.NewJaegerTracer(cfg, appConfig.App.AppName)
	if err != nil {
		panic(err)
	}
}

func initEtcd() {
	var err error
	cfg := &etcd.Config{}

	if err = resource.Config.ReadConfig("etcd", "toml", cfg); err != nil {
		panic(err)
	}

	resource.Etcd, err = etcd.NewClient(
		etcd.WithEndpoints(strings.Split(cfg.Endpoints, ";")),
		etcd.WithDialTimeout(cfg.DialTimeout),
	)
	if err != nil {
		panic(err)
	}
}

func initServices(ctx context.Context) {
	var (
		err   error
		dir   string
		files []string
	)

	if dir, err = filepath.Abs(confPath); err != nil {
		panic(err)
	}

	if files, err = filepath.Glob(filepath.Join(dir, "services", "*.toml")); err != nil {
		panic(err)
	}

	var discover registry.Discovery
	cfg := &servicer.Config{}
	for _, f := range files {
		f = path.Base(f)
		f = strings.TrimSuffix(f, path.Ext(f))

		if err = resource.Config.ReadConfig("services/"+f, "toml", cfg); err != nil {
			panic(err)
		}

		if cfg.Type == servicer.TypeRegistry {
			if resource.Etcd == nil {
				panic("initServices resource.Etcd nil")
			}
			opts := []registryEtcd.DiscoverOption{
				registryEtcd.WithContext(ctx),
				registryEtcd.WithServierName(cfg.ServiceName),
				registryEtcd.WithRefreshDuration(cfg.RefreshSecond),
				registryEtcd.WithDiscoverClient(resource.Etcd.Client),
			}
			if discover, err = registryEtcd.NewDiscovery(opts...); err != nil {
				panic(err)
			}
		}

		if err = servicer.LoadService(cfg, servicer.WithDiscovery(discover)); err != nil {
			panic(err)
		}
	}

	return
}

func initClientHTTP() {
	var err error
	cfg := &loggerRPC.RPCConfig{}

	if err = resource.Config.ReadConfig("log/rpc", "toml", cfg); err != nil {
		panic(err)
	}

	rpcLogger, err := loggerRPC.NewRPCLogger(cfg)
	if err != nil {
		panic(err)
	}

	resource.ClientHTTP = httpClient.New(
		httpClient.WithCodec(codec.JSONCodec{}),
		httpClient.WithLogger(rpcLogger),
		httpClient.WithBeforePlugins(&httpClient.JaegerBeforePlugin{}))
	if err != nil {
		panic(err)
	}
}
