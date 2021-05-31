package bootstrap

import (
	"gin-api/libraries/logging"
	"gin-api/libraries/mysql"
	"gin-api/libraries/redis"
	"gin-api/resource"
)

func Bootstrap() {
	initLogger()
	initMysql("default")
	initRedis("default")
}

func initMysql(db string) {
	var err error
	resource.TestDB, err = mysql.InitDB(db)
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
	resource.Logger = logging.NewLogger("./logs/gin-api.log", "./logs/gin-api.wf.log")
}

func initJaeger(err string) {
	return
}
