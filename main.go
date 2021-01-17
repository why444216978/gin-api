package main

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"gin-api/routers"
	"log"
	"strconv"
	"syscall"

	"github.com/why444216978/go-util/dir"

	"gin-api/libraries/endless"
)

func init() {
	logCfg := config.GetConfigToJson("log", "log")
	logDir := logCfg["dir"].(string) + "/" + app_const.SERVICE_NAME
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

func main() {
	server := routers.InitRouter()

	tmpServer := endless.NewServer(fmt.Sprintf(":%s", strconv.Itoa(app_const.SERVICE_PORT)), server)
	tmpServer.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := tmpServer.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
