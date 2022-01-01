package app

import (
	"context"
	"flag"
	"log"
	"path"
	"path/filepath"
	"strings"

	app_config "github.com/why444216978/gin-api/config"
	redis_cache "github.com/why444216978/gin-api/library/cache/redis"
	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/jaeger"
	jaeger_gorm "github.com/why444216978/gin-api/library/jaeger/gorm"
	jaeger_redis "github.com/why444216978/gin-api/library/jaeger/redis"
	redis_lock "github.com/why444216978/gin-api/library/lock/redis"
	"github.com/why444216978/gin-api/library/logging"
	logging_gorm "github.com/why444216978/gin-api/library/logging/gorm"
	logging_redis "github.com/why444216978/gin-api/library/logging/redis"
	logging_rpc "github.com/why444216978/gin-api/library/logging/rpc"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/redis"
	"github.com/why444216978/gin-api/library/registry"
	registry_etcd "github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/library/rpc/codec"
	"github.com/why444216978/gin-api/library/rpc/http"
	"github.com/why444216978/gin-api/library/servicer"
	"github.com/why444216978/gin-api/resource"
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
	flag.Parse()

	initConfig()
	initApp()
	initLogger()
	initMysql("test_mysql")
	initRedis("default_redis")
	initJaeger()
	initEtcd()
	initServices(ctx)
	initHTTPRPC()
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
	if err := resource.Config.ReadConfig("app", "toml", &app_config.App); err != nil {
		panic(err)
	}
}

func initLogger() {
	var err error
	cfg := &logging.Config{}

	if err = resource.Config.ReadConfig("log/service", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.ServiceLogger, err = logging.NewLogger(cfg,
		logging.WithModule(logging.ModuleHTTP),
		logging.WithServiceName(app_config.App.AppName),
	)
	if err != nil {
		panic(err)
	}

	RegisterCloseFunc(resource.ServiceLogger.Sync())
}

func initMysql(db string) {
	var err error
	cfg := &orm.Config{}
	logCfg := &logging_gorm.GormConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}

	if err = resource.Config.ReadConfig("log/gorm", "toml", logCfg); err != nil {
		panic(err)
	}

	logCfg.ServiceName = cfg.ServiceName
	gormLogger, err := logging_gorm.NewGorm(logCfg)
	if err != nil {
		panic(err)
	}

	resource.TestDB, err = orm.NewOrm(cfg,
		orm.WithTrace(jaeger_gorm.GormTrace),
		orm.WithLogger(gormLogger),
	)
	if err != nil {
		panic(err)
	}
}

func initRedis(db string) {
	var err error
	cfg := &redis.Config{}
	logCfg := &logging_redis.RedisConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}
	if err = resource.Config.ReadConfig("log/redis", "toml", &logCfg); err != nil {
		panic(err)
	}

	logCfg.ServiceName = cfg.ServiceName
	logCfg.Host = cfg.Host
	logCfg.Port = cfg.Port

	logger, err := logging_redis.NewRedisLogger(logCfg)
	if err != nil {
		panic(err)
	}

	rc := redis.NewClient(cfg)
	rc.AddHook(jaeger_redis.NewJaegerHook())
	rc.AddHook(logger)
	resource.RedisDefault = rc
}

func initLock() {
	var err error
	resource.RedisLock, err = redis_lock.New(resource.RedisDefault)
	if err != nil {
		panic(err)
	}
}

func initCache() {
	var err error

	resource.RedisCache, err = redis_cache.New(resource.RedisDefault, resource.RedisLock)
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

	_, _, err = jaeger.NewJaegerTracer(cfg, app_config.App.AppName)
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
			opts := []registry_etcd.DiscoverOption{
				registry_etcd.WithContext(ctx),
				registry_etcd.WithServierName(cfg.ServiceName),
				registry_etcd.WithRefreshDuration(cfg.RefreshSecond),
				registry_etcd.WithDiscoverClient(resource.Etcd.Client),
			}
			if discover, err = registry_etcd.NewDiscovery(opts...); err != nil {
				panic(err)
			}
		}

		if err = servicer.LoadService(cfg, servicer.WithDiscovery(discover)); err != nil {
			panic(err)
		}
	}

	return
}

func initHTTPRPC() {
	var err error
	cfg := &logging_rpc.RPCConfig{}

	if err = resource.Config.ReadConfig("log/rpc", "toml", cfg); err != nil {
		panic(err)
	}

	rpcLogger, err := logging_rpc.NewRPCLogger(cfg)
	if err != nil {
		panic(err)
	}

	resource.HTTPRPC = http.New(
		http.WithCodec(codec.JSONCodec{}),
		http.WithLogger(rpcLogger),
		http.WithBeforePlugins(&http.JaegerBeforePlugin{}))
	if err != nil {
		panic(err)
	}
}
