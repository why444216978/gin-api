package bootstrap

import (
	"flag"
	"log"
	"strings"

	"gin-api/global"
	"gin-api/libraries/config"
	"gin-api/libraries/etcd"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/libraries/orm"
	"gin-api/libraries/redis"
	"gin-api/resource"
)

var (
	conf = flag.String("conf", "conf_dev", "config path")
)

var confMap = map[string]struct{}{
	"conf_dev":      struct{}{},
	"conf_liantiao": struct{}{},
	"conf_qa":       struct{}{},
	"conf_online":   struct{}{},
}

func Bootstrap() {
	flag.Parse()

	initConfig()
	initApp()
	initLogger()
	initMysql("test_mysql")
	initRedis("default_redis")
	initJaeger()
	initEtcd()
}

func initConfig() {
	confPath := *conf
	log.Println("The conf path is :" + confPath)

	if _, ok := confMap[confPath]; !ok {
		panic(confPath + " error")
	}

	var err error
	resource.Config = config.InitConfig(confPath, "toml")
	if err != nil {
		panic(err)
	}
}

func initApp() {
	if err := resource.Config.ReadConfig("app", "toml", &global.Global); err != nil {
		panic(err)
	}
}

func initLogger() {
	var err error
	cfg := &logging.Config{}

	if err = resource.Config.ReadConfig("log", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.ServiceLogger, err = logging.NewLogger(cfg)
	if err != nil {
		panic(err)
	}
}

func initMysql(db string) {
	var err error
	cfg := &orm.Config{}
	gormCfg := &logging.GormConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}

	if err = resource.Config.ReadConfig("gorm_log", "toml", gormCfg); err != nil {
		panic(err)
	}

	gormLogger, err := logging.NewGorm(gormCfg)
	if err != nil {
		panic(err)
	}

	resource.TestDB, err = orm.NewOrm(cfg,
		orm.WithTrace(jaeger.GormTrace),
		orm.WithLogger(gormLogger),
	)
	if err != nil {
		panic(err)
	}
}

func initRedis(db string) {
	var err error
	cfg := &redis.Config{}
	logCfg := &logging.RedisConfig{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}
	if err = resource.Config.ReadConfig("redis_log", "toml", &logCfg); err != nil {
		panic(err)
	}

	logger, err := logging.NewRedisLogger(logCfg)
	if err != nil {
		panic(err)
	}

	rc := redis.NewClient(cfg)
	rc.AddHook(jaeger.NewJaegerHook())
	rc.AddHook(logger)
	resource.RedisCache = rc
}

func initJaeger() {
	var err error
	cfg := &jaeger.Config{}

	if err = resource.Config.ReadConfig("jaeger", "toml", cfg); err != nil {
		panic(err)
	}

	_, _, err = jaeger.NewJaegerTracer(cfg)
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
