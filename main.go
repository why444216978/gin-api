package main

import (
	"fmt"
	"gin-frame/configs"
	"log"
	"runtime"
	"strconv"
	"syscall"

	"gin-frame/routers"

	"gin-frame/libraries/endless"
)

var (
	port        int
	productName string
	moduleName  string
	env         string
	err         error
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	port = configs.SERVICE_PORT
	env = configs.ENV
	productName = configs.PRODUCT
	moduleName = configs.MODULE

	server := routers.InitRouter(port, productName, moduleName, env)

	tmpServer := endless.NewServer(fmt.Sprintf(":%s", strconv.Itoa(port)), server)
	tmpServer.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err = tmpServer.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
