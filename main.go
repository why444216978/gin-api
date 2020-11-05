package main

import (
	"fmt"
	"gin-api/configs"
	"gin-api/libraries/config"
	"gin-api/routers"
	"gin-api/libraries/logging"
	"log"
	"strconv"
	"syscall"

	"gin-api/libraries/endless"
)

func init() {
	logDir, _ := config.GetLogConfig(configs.LOG_SOURCE)
	file := configs.SERVICE_NAME + ".log"

	c := logging.LogConfig{
		Path:   logDir,
		File:   file,
		Mode:   1,
		Rotate: true,
		AsyncFormatter:true,
		RotatingFileHandler: logging.TIMED_ROTATING_FILE_HANDLER,
		RotateInterval: 3600,
		Debug:  true,
	}
	logging.Init(&c)
}

func main() {
	server := routers.InitRouter()

	tmpServer := endless.NewServer(fmt.Sprintf(":%s", strconv.Itoa(configs.SERVICE_PORT)), server)
	tmpServer.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := tmpServer.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
