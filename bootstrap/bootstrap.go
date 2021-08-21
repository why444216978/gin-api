package bootstrap

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"syscall"
	"time"

	"gin-api/global"
	"gin-api/jobs"
	"gin-api/libraries/config"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/libraries/orm"
	"gin-api/libraries/redis"
	"gin-api/resource"
	"gin-api/routers"
)

var (
	conf = flag.String("conf", "conf_dev", "config path")
	job  = flag.String("job", "", "is job")
)

func Bootstrap() {
	flag.Parse()

	initConfig()
	initApp()
	initLogger()
	initMysql("test_mysql")
	initRedis("default_redis")
	initJaeger()

	if *job == "" {
		log.Println("start by server")
		initHTTP()
	} else {
		jobs.Handle(*job)
	}
}

func initConfig() {
	log.Println("The conf path is :" + *conf)
	var err error
	resource.Config = config.InitConfig(*conf, "toml")
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
	var (
		err error
	)
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
	var (
		err error
	)
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
	var (
		err error
	)
	cfg := &redis.Config{}

	if err = resource.Config.ReadConfig(db, "toml", cfg); err != nil {
		panic(err)
	}

	rc := redis.NewClient(cfg)
	rc.AddHook(jaeger.NewJaegerHook())
	resource.RedisCache = rc
}

func initJaeger() {
	var (
		err error
	)
	cfg := &jaeger.Config{}

	if err = resource.Config.ReadConfig("jaeger", "toml", cfg); err != nil {
		panic(err)
	}

	_, _, err = jaeger.NewJaegerTracer(cfg)
	if err != nil {
		panic(err)
	}

	return
}

type HTTPConfig struct {
	Port int
}

func initHTTP() {
	router := routers.InitRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", global.Global.AppPort),
		Handler:      router,
		ReadTimeout:  time.Duration(global.Global.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(global.Global.WriteTimeout) * time.Millisecond,
	}
	log.Printf("Actual pid is %d", syscall.Getpid())
	log.Printf("Actual port is %d", global.Global.AppPort)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	// endless.DefaultReadTimeOut = 3 * time.Second
	// endless.DefaultWriteTimeOut = 3 * time.Second
	// serverEnd := endless.NewServer(fmt.Sprintf(":%d", global.Global.AppPort), router)
	// err = serverEnd.ListenAndServe()
	// if err != nil {
	// 	log.Printf("Server err: %v", err)
	// }
}
