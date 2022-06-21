package loader

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/why444216978/go-util/assert"
	utilDir "github.com/why444216978/go-util/dir"
	"github.com/why444216978/go-util/sys"

	"github.com/why444216978/gin-api/app/resource"
	httpClient "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/app"
	redisCache "github.com/why444216978/gin-api/library/cache/redis"
	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/jaeger"
	jaegerGorm "github.com/why444216978/gin-api/library/jaeger/gorm"
	jaegerRedis "github.com/why444216978/gin-api/library/jaeger/redis"
	redisLock "github.com/why444216978/gin-api/library/lock/redis"
	loggerGorm "github.com/why444216978/gin-api/library/logger/zap/gorm"
	loggerRedis "github.com/why444216978/gin-api/library/logger/zap/redis"
	loggerRPC "github.com/why444216978/gin-api/library/logger/zap/rpc"
	serviceLogger "github.com/why444216978/gin-api/library/logger/zap/service"
	"github.com/why444216978/gin-api/library/orm"
	"github.com/why444216978/gin-api/library/queue/rabbitmq"
	"github.com/why444216978/gin-api/library/redis"
	"github.com/why444216978/gin-api/library/registry"
	etcdRegistry "github.com/why444216978/gin-api/library/registry/etcd"
	registryEtcd "github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/library/servicer"
	"github.com/why444216978/gin-api/library/servicer/service"
	"github.com/why444216978/gin-api/server"
)

func Load() (err error) {
	// TODO
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	if err = loadLogger(); err != nil {
		return
	}
	if err = loadServices(); err != nil {
		return
	}
	if err = loadClientHTTP(); err != nil {
		return
	}
	// TODO 避免用户第一次使用运行panic，留给用户自己打开需要的依赖
	// if err = loadMysql("test_mysql"); err != nil {
	// 	return
	// }
	// if err = loadRedis("default_redis"); err != nil {
	// 	return
	// }
	// if err = loadJaeger(); err != nil {
	// 	return
	// }
	// if err = loadLock(); err != nil {
	// 	return
	// }
	// if err = loadCache(); err != nil {
	// 	return
	// }
	// if err = loadEtcd(); err != nil {
	// 	return
	// }
	// if err = loadRegistry(); err != nil {
	// 	return
	// }

	return
}

func loadLogger() (err error) {
	cfg := &serviceLogger.Config{}

	if err = config.ReadConfig("log/service", "toml", &cfg); err != nil {
		return
	}

	if resource.ServiceLogger, err = serviceLogger.NewServiceLogger(app.Name(), cfg); err != nil {
		return
	}

	server.RegisterCloseFunc(resource.ServiceLogger.Close())

	return
}

func loadMysql(db string) (err error) {
	cfg := &orm.Config{}
	logCfg := &loggerGorm.GormConfig{}

	if err = config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}

	if err = config.ReadConfig("log/gorm", "toml", logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	logger, err := loggerGorm.NewGorm(logCfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	if resource.TestDB, err = orm.NewOrm(cfg,
		orm.WithTrace(jaegerGorm.GormTrace),
		orm.WithLogger(logger),
	); err != nil {
		return
	}

	return
}

func loadRedis(db string) (err error) {
	cfg := &redis.Config{}
	logCfg := &loggerRedis.RedisConfig{}

	if err = config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}
	if err = config.ReadConfig("log/redis", "toml", &logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	logCfg.Host = cfg.Host
	logCfg.Port = cfg.Port

	logger, err := loggerRedis.NewRedisLogger(logCfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	rc := redis.NewClient(cfg)
	rc.AddHook(jaegerRedis.NewJaegerHook())
	rc.AddHook(logger)
	resource.RedisDefault = rc

	return
}

func loadRabbitMQ(service string) (err error) {
	cfg := &rabbitmq.Config{}
	if err = config.ReadConfig(service, "toml", cfg); err != nil {
		return
	}

	if resource.RabbitMQ, err = rabbitmq.New(cfg); err != nil {
		return
	}

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

	if err = config.ReadConfig("jaeger", "toml", cfg); err != nil {
		return
	}

	if _, _, err = jaeger.NewJaegerTracer(cfg, app.Name()); err != nil {
		return
	}

	return
}

func loadEtcd() (err error) {
	cfg := &etcd.Config{}

	if err = config.ReadConfig("etcd", "toml", cfg); err != nil {
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

	if err = config.ReadConfig("registry", "toml", cfg); err != nil {
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
		etcdRegistry.WithRegistrarServiceName(app.Name()),
		etcdRegistry.WithRegistarHost(localIP),
		etcdRegistry.WithRegistarPort(app.Port()),
		etcdRegistry.WithRegistrarLease(cfg.Lease)); err != nil {
		return
	}

	if err = server.RegisterCloseFunc(resource.Registrar.DeRegister); err != nil {
		return
	}

	return
}

func loadServices() (err error) {
	var (
		dir   string
		files []string
	)

	if dir, err = filepath.Abs(config.Path()); err != nil {
		return
	}

	if files, err = filepath.Glob(filepath.Join(dir, "services", "*.toml")); err != nil {
		return
	}

	var discover registry.Discovery
	info := utilDir.FileInfo{}
	cfg := &service.Config{}
	for _, f := range files {
		if info, err = utilDir.GetPathInfo(f); err != nil {
			return
		}
		if err = config.ReadConfig("services/"+info.BaseNoExt, info.ExtNoSpot, cfg); err != nil {
			return
		}

		if cfg.Type == servicer.TypeRegistry {
			if assert.IsNil(resource.Etcd) {
				return errors.New("loadServices resource.Etcd nil")
			}
			opts := []registryEtcd.DiscoverOption{
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
	if err = config.ReadConfig("log/rpc", "toml", cfg); err != nil {
		return
	}

	logger, err := loggerRPC.NewRPCLogger(cfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	resource.ClientHTTP = httpClient.New(
		httpClient.WithLogger(logger),
		httpClient.WithBeforePlugins(&httpClient.JaegerBeforePlugin{}))
	if err != nil {
		return
	}

	return
}
