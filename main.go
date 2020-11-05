package main

import (
	"fmt"
	"gin-api/configs"
	"log"
	"runtime"
	"strconv"
	"syscall"

	"gin-api/routers"

	"gin-api/libraries/endless"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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
