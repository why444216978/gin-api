package bootstrap

import (
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
	"gin-api/resource"

	"github.com/why444216978/go-util/dir"
)

func Bootstrap() {
	initLogger()
	initMysql("default")
	initRedis("default")
}

func initMysql(db string) {
	var err error
	resource.TestDB, err = mysql.GetConn(db)
	if err != nil {
		panic(err)
	}
}

func initRedis(db string) {
	var err error
	resource.DefaultRedis, err = redis.GetRedis(db)
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logCfg := config.GetConfigToJson("log", "log")
	logDir := logCfg["dir"].(string)
	file := app_const.SERVICE_NAME + ".log"
	dir.CreateDir(logDir)
	c := logging.LogConfig{
		Path:                logDir,
		File:                file,
		Mode:                1,
		Rotate:              true,
		AsyncFormatter:      true,
		RotatingFileHandler: logging.TIMED_ROTATING_FILE_HANDLER,
		RotateInterval:      3600,
		Debug:               true,
	}
	logging.Init(&c)
}

func initJaeger(db string) {
	var err error
	resource.TestDB, err = mysql.GetConn(db)
	if err != nil {
		panic(err)
	}
}
