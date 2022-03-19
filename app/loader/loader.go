package loader

import (
	"context"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/why444216978/go-util/assert"
	"github.com/why444216978/go-util/sys"

	appConfig "github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/resource"
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
	etcdRegistry "github.com/why444216978/gin-api/library/registry/etcd"
	registryEtcd "github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/library/servicer"
	"github.com/why444216978/gin-api/library/servicer/service"
	"github.com/why444216978/gin-api/server"
)

var envFlag = flag.String("env", "dev", "config path")

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

func Load() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	if err = loadConfig(); err != nil {
		return
	}
	if err = loadApp(); err != nil {
		return
	}
	if err = loadLogger(); err != nil {
		return
	}
	if err = loadServices(ctx); err != nil {
		return
	}
	if err = loadClientHTTP(); err != nil {
		return
	}
	// TODO 避免用户第一次使用运行panic，留给用户自己打开需要的依赖
	if err = loadMysql("test_mysql"); err != nil {
		return
	}
	if err = loadRedis("default_redis"); err != nil {
		return
	}
	if err = loadJaeger(); err != nil {
		return
	}
	if err = loadLock(); err != nil {
		return
	}
	if err = loadCache(); err != nil {
		return
	}
	if err = loadEtcd(); err != nil {
		return
	}
	if err = loadRegistry(); err != nil {
		return
	}

	return
}

func loadConfig() (err error) {
	env = *envFlag
	log.Println("The environment is :" + env)

	if _, ok := envMap[env]; !ok {
		panic(env + " error")
	}

	confPath = "conf/" + env
	if _, err = os.Stat(confPath); err != nil {
		return
	}

	resource.Config = config.InitConfig(confPath, "toml")

	return
}

func loadApp() (err error) {
	return resource.Config.ReadConfig("app", "toml", &appConfig.App)
}

func loadLogger() (err error) {
	cfg := &logger.Config{}

	if err = resource.Config.ReadConfig("log/service", "toml", &cfg); err != nil {
		return
	}

	if resource.ServiceLogger, err = logger.NewLogger(cfg,
		logger.WithModule(logger.ModuleHTTP),
		logger.WithServiceName(appConfig.App.AppName),
	); err != nil {
		return
	}

	server.RegisterCloseFunc(resource.ServiceLogger.Sync())

	return
}

func loadMysql(db string) (err error) {
	cfg := &orm.Config{}
	logCfg := &loggerGorm.GormConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}

	if err = resource.Config.ReadConfig("log/gorm", "toml", logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	gormLogger, err := loggerGorm.NewGorm(logCfg)
	if err != nil {
		return
	}

	if resource.TestDB, err = orm.NewOrm(cfg,
		orm.WithTrace(jaegerGorm.GormTrace),
		orm.WithLogger(gormLogger),
	); err != nil {
		return
	}

	return
}

func loadRedis(db string) (err error) {
	cfg := &redis.Config{}
	logCfg := &loggerRedis.RedisConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}
	if err = resource.Config.ReadConfig("log/redis", "toml", &logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	logCfg.Host = cfg.Host
	logCfg.Port = cfg.Port

	logger, err := loggerRedis.NewRedisLogger(logCfg)
	if err != nil {
		return
	}

	rc := redis.NewClient(cfg)
	rc.AddHook(jaegerRedis.NewJaegerHook())
	rc.AddHook(logger)
	resource.RedisDefault = rc

	return
}

func loadLock() (err error) {
	resource.RedisLock, err = redisLock.New(resource.RedisDefault)
	return
}

func loadCache() (err error) {
	resource.RedisCache, err = redisCache.New(resource.RedisDefault, resource.RedisLock)
	return
}

func loadJaeger() (err error) {
	cfg := &jaeger.Config{}

	if err = resource.Config.ReadConfig("jaeger", "toml", cfg); err != nil {
		return
	}

	if _, _, err = jaeger.NewJaegerTracer(cfg, appConfig.App.AppName); err != nil {
		return
	}

	return
}

func loadEtcd() (err error) {
	cfg := &etcd.Config{}

	if err = resource.Config.ReadConfig("etcd", "toml", cfg); err != nil {
		return
	}

	if resource.Etcd, err = etcd.NewClient(
		etcd.WithEndpoints(strings.Split(cfg.Endpoints, ";")),
		etcd.WithDialTimeout(cfg.DialTimeout),
	); err != nil {
		return
	}

	return
}

func loadRegistry() (err error) {
	var (
		localIP string
		cfg     = &registry.RegistryConfig{}
	)

	if err = resource.Config.ReadConfig("registry", "toml", cfg); err != nil {
		return
	}

	if localIP, err = sys.LocalIP(); err != nil {
		return
	}

	if assert.IsNil(resource.Etcd) {
		err = errors.New("resource.Etcd is nil")
		return
	}

	if resource.Registrar, err = etcdRegistry.NewRegistry(
		etcdRegistry.WithRegistrarClient(resource.Etcd.Client),
		etcdRegistry.WithRegistrarServiceName(appConfig.App.AppName),
		etcdRegistry.WithRegistarHost(localIP),
		etcdRegistry.WithRegistarPort(appConfig.App.AppPort),
		etcdRegistry.WithRegistrarLease(cfg.Lease)); err != nil {
		return
	}

	if err = server.RegisterCloseFunc(resource.Registrar.DeRegister); err != nil {
		return
	}

	return
}

func loadServices(ctx context.Context) (err error) {
	var (
		dir   string
		files []string
	)

	if dir, err = filepath.Abs(confPath); err != nil {
		return
	}

	if files, err = filepath.Glob(filepath.Join(dir, "services", "*.toml")); err != nil {
		return
	}

	var discover registry.Discovery
	cfg := &service.Config{}
	for _, f := range files {
		f = path.Base(f)
		f = strings.TrimSuffix(f, path.Ext(f))

		if err = resource.Config.ReadConfig("services/"+f, "toml", cfg); err != nil {
			return
		}

		if cfg.Type == servicer.TypeRegistry {
			if assert.IsNil(resource.Etcd) {
				return errors.New("loadServices resource.Etcd nil")
			}
			opts := []registryEtcd.DiscoverOption{
				registryEtcd.WithContext(ctx),
				registryEtcd.WithServierName(cfg.ServiceName),
				registryEtcd.WithRefreshDuration(cfg.RefreshSecond),
				registryEtcd.WithDiscoverClient(resource.Etcd.Client),
			}
			if discover, err = registryEtcd.NewDiscovery(opts...); err != nil {
				return
			}
		}

		if err = service.LoadService(cfg, service.WithDiscovery(discover)); err != nil {
			return
		}
	}

	return
}

func loadClientHTTP() (err error) {
	cfg := &loggerRPC.RPCConfig{}
	if err = resource.Config.ReadConfig("log/rpc", "toml", cfg); err != nil {
		return
	}

	rpcLogger, err := loggerRPC.NewRPCLogger(cfg)
	if err != nil {
		return
	}

	resource.ClientHTTP = httpClient.New(
		httpClient.WithLogger(rpcLogger),
		httpClient.WithBeforePlugins(&httpClient.JaegerBeforePlugin{}))
	if err != nil {
		return
	}

	return
}
