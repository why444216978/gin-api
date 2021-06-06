package bootstrap

import (
	"gin-api/libraries/config"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
	"gin-api/resource"
)

func Bootstrap() {
	initConfig()
	initLogger()
	initMysql("default")
	initRedis("default")
	initJaeger()
}

func initConfig() {
	var err error
	resource.Config = config.InitConfig("./conf_dev", "toml")
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	var (
		err error
		cfg logging.Config
	)

	if err = resource.Config.ReadConfig("log", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.Logger = logging.NewLogger(cfg)
}

func initMysql(db string) {
	var (
		err error
		cfg mysql.Config
	)

	if err = resource.Config.ReadConfig("test_mysql", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.TestDB, err = mysql.NewMySQL(cfg)
	if err != nil {
		panic(err)
	}
}

func initRedis(db string) {
	var (
		err error
		cfg redis.Config
	)

	if err = resource.Config.ReadConfig("default_redis", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.DefaultRedis, err = redis.GetRedis(cfg)
	if err != nil {
		panic(err)
	}
}

func initJaeger() {
	var (
		err error
		cfg jaeger.Config
	)

	if err = resource.Config.ReadConfig("jaeger", "toml", &cfg); err != nil {
		panic(err)
	}

	resource.Tracer, _, err = jaeger.NewJaegerTracer(cfg)
	if err != nil {
		panic(err)
	}

	return
}
